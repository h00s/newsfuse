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

func NewKliknihr(cb func(h models.Headline)) internal.Scraper {
	s := internal.NewScraper("kliknihr", "https://www.klikni.hr", 15, 75)
	s.RegisterOnHeadline(cb)
	s.Collector.OnHTML("h3.jeg_post_title", func(e *colly.HTMLElement) {
		s.HeadlineCallback(models.Headline{
			Title:       e.ChildText("a"),
			URL:         e.ChildAttr("a", "href"),
			PublishedAt: time.Now(),
		})
	})
	s.Start()
	return s
}
