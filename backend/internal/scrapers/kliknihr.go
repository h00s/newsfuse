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
	s.Collector.OnHTML("h3.jeg_post_title", func(e *colly.HTMLElement) {
		s.Headline <- models.Headline{
			Source:      s.Name,
			Title:       e.ChildText("a"),
			URL:         e.ChildAttr("a", "href"),
			PublishedAt: time.Now(),
		}
	})
	s.Start()
	return s
}
