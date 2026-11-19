package sites

import (
	"context"

	"github.com/h00s/newsfuse/internal/collector"
	"github.com/h00s/newsfuse/internal/scraper"
)

type N1InfoCroatia struct {
	scraper.Base
}

func NewN1InfoCroatia() *N1InfoCroatia {
	s := &N1InfoCroatia{
		Base: *scraper.New(
			"N1",
			"https://n1info.hr/vijesti/",
			10,
			20,
			[]int{0, 1, 2, 3, 4, 5, 23},
		),
	}

	s.ScrapeHeadline("h3[data-testid='article-title']", func(e *collector.HTMLElement) {
		s.AddHeadline(scraper.Headline{
			Title: e.ChildText("a"),
			URL:   "https://n1info.hr" + e.ChildAttr("a", "href"),
		})
	})

	return s
}

func (s *N1InfoCroatia) ScrapeStory(ctx context.Context, url string) (string, error) {
	return s.ScrapeStoryFrom(ctx, url, "div.article-content-wrapper", "p[data-block-key]", false)
}
