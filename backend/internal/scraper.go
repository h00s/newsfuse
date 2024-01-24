package internal

import (
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
	MinRefreshInterval int
	MaxRefreshInterval int
	HeadlineCallback   func(h models.Headlines)
	Collector          *colly.Collector
}

func NewScraper(name, url string, minRefreshInterval, maxRefreshInterval int) *DefaultScraper {
	return &DefaultScraper{
		URL:                url,
		MinRefreshInterval: minRefreshInterval,
		MaxRefreshInterval: maxRefreshInterval,
		Collector:          colly.NewCollector(),
	}
}

/*func (s *DefaultScraper) OnHTML(selector string, cb func(e *colly.HTMLElement)) {
	s.Collector.OnHTML(selector, cb)
}*/

func (s *DefaultScraper) OnHeadline(cb func(h models.Headline)) {
}

func (s *DefaultScraper) Start() {
	go func() {
		for {
			s.Collector.Visit(s.URL)
			s.Collector.Wait()
			waitTime := rand.Intn(s.MaxRefreshInterval-s.MinRefreshInterval) + s.MinRefreshInterval
			time.Sleep(time.Duration(waitTime) * time.Minute)
		}
	}()
}
