package internal

import (
	"fmt"
	"log/slog"
	"math/rand"
	"slices"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
)

type Scraper interface {
	Start()
	ScrapeStory(url string) (string, error)
}

type DefaultScraper struct {
	headlinesChannel chan models.Headlines
	headlines        models.Headlines

	Name string
	URL  string

	MinRefreshInterval int
	MaxRefreshInterval int

	OffHours []int

	log       *slog.Logger
	collector *colly.Collector
}

func NewScraper(headlinesChannel chan models.Headlines, log *slog.Logger, name, url string, minRefreshInterval, maxRefreshInterval int, offHours []int) *DefaultScraper {
	return &DefaultScraper{
		headlinesChannel: headlinesChannel,
		headlines:        nil,

		Name: name,
		URL:  url,

		MinRefreshInterval: minRefreshInterval,
		MaxRefreshInterval: maxRefreshInterval,

		OffHours: offHours,

		log:       log,
		collector: colly.NewCollector(),
	}
}

func (s *DefaultScraper) AddHeadline(h models.Headline) {
	title := strings.TrimSpace(h.Title)
	if title != "" {
		h.Title = title
		h.URL = strings.TrimSpace(h.URL)
		s.headlines = append(s.headlines, h)
	}
}

func (s *DefaultScraper) ScrapeHeadline(selector string, callback func(e *colly.HTMLElement)) {
	s.collector.OnHTML(selector, callback)
}

func (s *DefaultScraper) ScrapeStory(url, element, childElement string, html bool) (string, error) {
	var story string

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"

	c.OnHTML(element, func(e *colly.HTMLElement) {
		var contents string
		e.ForEach(childElement, func(_ int, el *colly.HTMLElement) {
			if html {
				contentsHTML, err := el.DOM.Html()
				if err != nil {
					s.log.Error("Error getting HTML", "error", err.Error())
					contents = ""
				} else {
					contents = strings.TrimSpace(contentsHTML)
				}
			} else {
				contents = strings.TrimSpace(el.Text)
			}
			story += fmt.Sprintf("<p>%s</p>", contents)
		})
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	return story, nil
}

func (s *DefaultScraper) Start() {
	s.collector.DisableCookies()
	s.collector.AllowURLRevisit = true

	s.collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
		s.headlines = models.Headlines{}
	})

	s.collector.OnScraped(func(r *colly.Response) {
		s.headlinesChannel <- s.headlines
		s.log.Info("Finished scraping", "scraper", s.Name)
	})

	go func() {
		for {
			s.log.Info("Started scraping", "scraper", s.Name)
			if !slices.Contains(s.OffHours, time.Now().Hour()) {
				s.collector.Visit(s.URL)
				s.collector.Wait()
			} else {
				s.log.Info("Skipping scraping", "scraper", s.Name, "hour", time.Now().Hour())
			}
			waitTime := rand.Intn(s.MaxRefreshInterval-s.MinRefreshInterval) + s.MinRefreshInterval
			s.log.Info("Waiting for next scraping", "scraper", s.Name, "minutes", waitTime)
			time.Sleep(time.Duration(waitTime) * time.Minute)
		}
	}()
}
