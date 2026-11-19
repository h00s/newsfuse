// Package sites implements the scraper for each news source.
package sites

import (
	"context"
	"strings"

	"github.com/h00s/newsfuse/internal/collector"
	"github.com/h00s/newsfuse/internal/scraper"
)

type Bughr struct {
	scraper.Base
}

func NewBughr() *Bughr {
	s := &Bughr{
		Base: *scraper.New(
			"Bug",
			"https://www.bug.hr",
			10,
			30,
			[]int{0, 1, 2, 3, 4, 5, 23},
		),
	}

	s.ScrapeHeadline("h2[class='post-listing__title'] > a", func(e *collector.HTMLElement) {
		// skip premium articles
		if e.DOM.Find("svg").Length() > 0 {
			return
		}

		url := e.Attr("href")
		// skip forum articles
		if strings.HasPrefix(url, "/forum") {
			return
		}
		if !strings.HasPrefix(url, "http") {
			url = s.URL + url
		}

		s.AddHeadline(scraper.Headline{
			Title: e.Text,
			URL:   url,
		})
	})

	return s
}

func (s *Bughr) ScrapeStory(ctx context.Context, url string) (string, error) {
	return s.ScrapeStoryFrom(ctx, url, "div[class^='post-full__content']", "p:not([class])", false)
}
