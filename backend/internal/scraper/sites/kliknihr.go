package sites

import (
	"context"

	"github.com/h00s/newsfuse/internal/collector"
	"github.com/h00s/newsfuse/internal/scraper"
)

type Kliknihr struct {
	scraper.Base
}

func NewKliknihr() *Kliknihr {
	s := &Kliknihr{
		Base: *scraper.New(
			"klikni.hr",
			"https://www.klikni.hr",
			15,
			30,
			[]int{0, 1, 2, 3, 4, 5, 23},
		),
	}

	s.ScrapeHeadline("h3.jeg_post_title", func(e *collector.HTMLElement) {
		s.AddHeadline(scraper.Headline{
			Title: e.ChildText("a"),
			URL:   e.ChildAttr("a", "href"),
		})
	})

	return s
}

func (s *Kliknihr) ScrapeStory(ctx context.Context, url string) (string, error) {
	return s.ScrapeStoryFrom(ctx, url, "section[class='container page-content']", "p:not([class])", false)
}
