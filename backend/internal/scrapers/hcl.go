package scrapers

import (
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
)

type HCL struct {
	internal.DefaultScraper
}

func NewHCL(h chan (models.Headlines), sourceID int64) *HCL {
	s := &HCL{
		DefaultScraper: *internal.NewScraper(h,
			"HCL",
			"https://www.hcl.hr/",
			10,
			20,
			[]int{0, 1, 2, 3, 4, 5, 23},
		),
	}

	s.ScrapeHeadline("div[id^='post-']", func(e *colly.HTMLElement) {
		title := ""
		e.ForEach("div.text h2 a", func(_ int, el *colly.HTMLElement) {
			sel := el.DOM
			sel.Find("span.title").Remove()
			title = strings.TrimSpace(sel.Text())
		})

		s.AddHeadline(models.Headline{
			SourceID:    sourceID,
			Title:       title,
			URL:         e.ChildAttr("div.text h2 a", "href"),
			PublishedAt: time.Now(),
		})
	})

	return s
}

func (s *HCL) ScrapeStory(url string) (string, error) {
	return s.DefaultScraper.ScrapeStory(url, "div.article", "p:not([class]):not(.meta p, .tags p)", false)
}
