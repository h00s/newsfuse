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

func NewKliknihr(h chan (models.Headline)) internal.Scraper {
	s := internal.NewScraper("klikni.hr", "https://www.klikni.hr", 15, 30, h)

	s.ScrapeHeadline("h3.jeg_post_title", func(e *colly.HTMLElement) {
		s.AddHeadline(models.Headline{
			Source:      s.Name,
			Title:       e.ChildText("a"),
			URL:         e.ChildAttr("a", "href"),
			PublishedAt: time.Now(),
		})
	})

	return s
}
