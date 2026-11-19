package sites

import (
	"context"
	"strings"

	"github.com/h00s/newsfuse/internal/collector"
	"github.com/h00s/newsfuse/internal/scraper"
)

type HCL struct {
	scraper.Base
}

func NewHCL() *HCL {
	s := &HCL{
		Base: *scraper.New(
			"HCL",
			"https://www.hcl.hr/",
			10,
			20,
			[]int{0, 1, 2, 3, 4, 5, 23},
		),
	}

	s.ScrapeHeadline("div.articles div[id^='post-']", func(e *collector.HTMLElement) {
		var title string
		e.ForEach("div.text h2 a", func(_ int, el *collector.HTMLElement) {
			sel := el.DOM.Clone()
			sel.Find("span.title").Remove()
			title = strings.TrimSpace(sel.Text())
		})

		url := e.ChildAttr("div.text h2 a", "href")
		if title == "" || url == "" {
			return
		}

		s.AddHeadline(scraper.Headline{
			Title: title,
			URL:   url,
		})
	})

	return s
}

func (s *HCL) ScrapeStory(ctx context.Context, url string) (string, error) {
	return s.ScrapeStoryFrom(ctx, url, "div.article", "p:not([class]):not(.meta p, .tags p)", false)
}
