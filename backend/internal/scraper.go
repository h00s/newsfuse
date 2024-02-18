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

	headlineChannel chan (models.Headline)
	headlines       models.Headlines

	utils     *raptor.Utils
	collector *colly.Collector
}

func NewScraper(name, url string, minRefreshInterval, maxRefreshInterval int, headlineChannel chan (models.Headline)) *DefaultScraper {
	return &DefaultScraper{
		Name: name,
		URL:  url,

		MinRefreshInterval: minRefreshInterval,
		MaxRefreshInterval: maxRefreshInterval,

		headlineChannel: headlineChannel,
		headlines:       nil,

		utils:     nil,
		collector: colly.NewCollector(),
	}
}

func (s *DefaultScraper) Init(u *raptor.Utils) {
	s.utils = u
}

func (s *DefaultScraper) AddHeadline(h models.Headline) {
	if strings.TrimSpace(h.Title) != "" {
		s.headlines = append(s.headlines, h)
	}
}

func (s *DefaultScraper) ScrapeHeadline(selector string, callback func(e *colly.HTMLElement)) {
	s.collector.OnHTML(selector, callback)
}

func (s *DefaultScraper) ScrapeStory(url, element, childElement string) (string, error) {
	var story string

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"

	c.OnHTML(element, func(e *colly.HTMLElement) {
		e.ForEach(childElement, func(_ int, el *colly.HTMLElement) {
			story += fmt.Sprintf("<p>%s</p>", strings.TrimSpace(el.Text))
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
		for i := len(s.headlines) - 1; i >= 0; i-- {
			h := s.headlines[i]
			s.headlineChannel <- h
		}
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
