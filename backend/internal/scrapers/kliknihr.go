package scrapers

import (
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
)

type Kliknihr struct {
	internal.DefaultScraper
}

func NewKliknihr(headline chan (models.Headline)) internal.Scraper {
	s := internal.NewScraper(headline, "klikni.hr", "https://www.klikni.hr", 15, 75)
	s.Collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
		s.Headlines = models.Headlines{}
	})
	s.Collector.OnHTML("h3.jeg_post_title", func(e *colly.HTMLElement) {
		headline := models.Headline{
			Source:      s.Name,
			Title:       e.ChildText("a"),
			URL:         e.ChildAttr("a", "href"),
			PublishedAt: time.Now(),
		}
		s.Headlines = append(s.Headlines, headline)
	})
	s.Collector.OnScraped(func(r *colly.Response) {
		for i := len(s.Headlines) - 1; i >= 0; i-- {
			h := s.Headlines[i]
			s.Headline <- h
		}
	})
	s.Start()
	return s
}
