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
	Headline           chan (models.Headline)
	Collector          *colly.Collector
	Headlines          models.Headlines
}

func NewScraper(headline chan (models.Headline), name, url string, minRefreshInterval, maxRefreshInterval int) *DefaultScraper {
	return &DefaultScraper{
		Headline:           headline,
		Name:               name,
		URL:                url,
		MinRefreshInterval: minRefreshInterval,
		MaxRefreshInterval: maxRefreshInterval,
		Collector:          colly.NewCollector(),
	}
}

func (s *DefaultScraper) Start() {
	s.Collector.DisableCookies()
	s.Collector.AllowURLRevisit = true
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
