package scrapers

import (
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
)

type Mojportalhr struct {
	internal.DefaultScraper
	Source models.Source
}

func NewMojportalhr(h chan (models.Headlines), sourceID int64) *Mojportalhr {
	s := &Mojportalhr{
		DefaultScraper: *internal.NewScraper(h,
			"MojPortal.hr",
			"https://www.mojportal.hr/",
			15,
			30,
			[]int{0, 1, 2, 3, 4, 5, 23},
		),
	}

	s.ScrapeHeadline("a.se-card--link", func(e *colly.HTMLElement) {
		s.AddHeadline(models.Headline{
			SourceID:    sourceID,
			Title:       strings.TrimSpace(e.ChildText("div.se-card--head")),
			URL:         "https://www.mojportal.hr" + strings.TrimSpace(e.Attr("href")),
			PublishedAt: time.Now(),
		})
	})

	return s
}

func (s *Mojportalhr) ScrapeStory(url string) (string, error) {
	return s.DefaultScraper.ScrapeStory(url, "div.se-article--text", "p:not([class])", false)
}
