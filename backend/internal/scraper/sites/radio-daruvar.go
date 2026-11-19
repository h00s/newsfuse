package sites

import (
	"context"
	"strings"

	"github.com/h00s/newsfuse/internal/collector"
	"github.com/h00s/newsfuse/internal/scraper"
)

type RadioDaruvar struct {
	scraper.Base
}

func NewRadioDaruvar() *RadioDaruvar {
	s := &RadioDaruvar{
		Base: *scraper.New(
			"Radio Daruvar",
			"https://www.radio-daruvar.hr/",
			30,
			60,
			[]int{0, 1, 2, 3, 4, 5, 19, 20, 21, 22, 23},
		),
	}

	s.ScrapeHeadline("li[class='news-item'],h2[class='post-title']", func(e *collector.HTMLElement) {
		s.AddHeadline(scraper.Headline{
			Title: strings.TrimSpace(e.ChildText("a")),
			URL:   strings.TrimSpace(e.ChildAttr("a", "href")),
		})
	})

	return s
}

func (s *RadioDaruvar) ScrapeStory(ctx context.Context, url string) (string, error) {
	return s.ScrapeStoryFrom(ctx, url, "div[class='entry-content entry clearfix']", "p:not([class])", false)
}
