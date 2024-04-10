package internal

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-raptor/raptor"
	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
)

type Scraper interface {
	Start()
	Init(u *raptor.Utils)
	ScrapeStory(url string) (string, error)
}

type DefaultScraper struct {
	Name string
	URL  string

	MinRefreshInterval int
	MaxRefreshInterval int

	headlinesChannel chan (models.Headlines)
	headlines        models.Headlines

	utils     *raptor.Utils
	collector *colly.Collector
}

func NewScraper(name, url string, minRefreshInterval, maxRefreshInterval int, headlinesChannel chan (models.Headlines)) *DefaultScraper {
	return &DefaultScraper{
		Name: name,
		URL:  url,

		MinRefreshInterval: minRefreshInterval,
		MaxRefreshInterval: maxRefreshInterval,

		headlinesChannel: headlinesChannel,
		headlines:        nil,

		utils:     nil,
		collector: colly.NewCollector(),
	}
}

func (s *DefaultScraper) Init(u *raptor.Utils) {
	s.utils = u
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
					s.utils.Log.Error("Error getting HTML", "error", err.Error())
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
		s.utils.Log.Info("Finished scraping", "scraper", s.Name)
	})

	go func() {
		for {
			s.utils.Log.Info("Started scraping", "scraper", s.Name)
			s.collector.Visit(s.URL)
			s.collector.Wait()
			waitTime := rand.Intn(s.MaxRefreshInterval-s.MinRefreshInterval) + s.MinRefreshInterval
			s.utils.Log.Info("Waiting for next scraping", "scraper", s.Name, "minutes", waitTime)
			time.Sleep(time.Duration(waitTime) * time.Minute)
		}
	}()
}
