package scrapers

import (
	"context"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/app/utils"
)

type HCL struct {
	utils.DefaultScraper
}

func NewHCL(sourceID int64) *HCL {
	s := &HCL{
		DefaultScraper: *utils.NewScraper(
			"HCL",
			"https://www.hcl.hr/",
			10,
			20,
			[]int{0, 1, 2, 3, 4, 5, 23},
		),
	}

	s.ScrapeHeadline("div.articles div[id^='post-']", func(e *colly.HTMLElement) {
		var title string
		e.ForEach("div.text h2 a", func(_ int, el *colly.HTMLElement) {
			sel := el.DOM.Clone()
			sel.Find("span.title").Remove()
			title = strings.TrimSpace(sel.Text())
		})

		url := e.ChildAttr("div.text h2 a", "href")
		if title == "" || url == "" {
			return
		}

		s.AddHeadline(models.Headline{
			SourceID:    sourceID,
			Title:       title,
			URL:         url,
			PublishedAt: time.Now(),
		})
	})

	return s
}

func (s *HCL) ScrapeStory(ctx context.Context, url string) (string, error) {
	return s.ScrapeStoryFrom(ctx, url, "div.article", "p:not([class]):not(.meta p, .tags p)", false)
}
