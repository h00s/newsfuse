package internal

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
)

type Scraper interface {
	Start()
}

type DefaultScraper struct {
	URL                string
	Name               string
	MinRefreshInterval int
	MaxRefreshInterval int
	HeadlineChannel    chan (models.Headline)
	Collector          *colly.Collector
	headlines          models.Headlines
}

func NewScraper(headline chan (models.Headline), name, url string, minRefreshInterval, maxRefreshInterval int) *DefaultScraper {
	return &DefaultScraper{
		HeadlineChannel:    headline,
		Name:               name,
		URL:                url,
		MinRefreshInterval: minRefreshInterval,
		MaxRefreshInterval: maxRefreshInterval,
		Collector:          colly.NewCollector(),
	}
}

func (s *DefaultScraper) AddHeadline(h models.Headline) {
	s.headlines = append(s.headlines, h)
}

func (s *DefaultScraper) Start() {
	s.Collector.DisableCookies()
	s.Collector.AllowURLRevisit = true

	s.Collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
		s.headlines = models.Headlines{}
	})

	s.Collector.OnScraped(func(r *colly.Response) {
		for i := len(s.headlines) - 1; i >= 0; i-- {
			h := s.headlines[i]
			s.HeadlineChannel <- h
		}
	})

	go func() {
		for {
			fmt.Println("Scraping", s.URL)
			s.Collector.Visit(s.URL)
			s.Collector.Wait()
			waitTime := rand.Intn(s.MaxRefreshInterval-s.MinRefreshInterval) + s.MinRefreshInterval
			time.Sleep(time.Duration(waitTime) * time.Minute)
		}
	}()
}
