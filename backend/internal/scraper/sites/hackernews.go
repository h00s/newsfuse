package sites

import (
	"context"

	"github.com/h00s/newsfuse/internal/collector"
	"github.com/h00s/newsfuse/internal/scraper"
)

type HackerNews struct {
	scraper.Base
}

func NewHackerNews() *HackerNews {
	s := &HackerNews{
		Base: *scraper.New(
			"Hacker News",
			"https://news.ycombinator.com/",
			10,
			15,
			[]int{},
		),
	}

	s.ScrapeHeadline("tr[class^='athing']", func(e *collector.HTMLElement) {
		url := s.URL + "item?id=" + e.Attr("id")
		anchor := e.DOM.Find("span[class='titleline'] > a").First()
		s.AddHeadline(scraper.Headline{
			Title: anchor.Text(),
			URL:   url,
		})
	})

	return s
}

func (s *HackerNews) ScrapeStory(ctx context.Context, url string) (string, error) {
	return s.ScrapeStoryFrom(ctx, url, "td[class='title']", "span[class='titleline']", true)
}
