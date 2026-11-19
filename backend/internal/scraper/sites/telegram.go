package sites

import (
	"context"

	"github.com/h00s/newsfuse/internal/collector"
	"github.com/h00s/newsfuse/internal/scraper"
)

type Telegram struct {
	scraper.Base
}

func NewTelegram() *Telegram {
	s := &Telegram{
		Base: *scraper.New(
			"Telegram",
			"https://www.telegram.hr/vijesti",
			15,
			30,
			[]int{1, 2, 3, 4, 5},
		),
	}

	s.ScrapeHeadline("a[role='article']", func(e *collector.HTMLElement) {
		s.AddHeadline(scraper.Headline{
			Title: e.ChildText("h2"),
			URL:   "https://www.telegram.hr" + e.Attr("href"),
		})
	})

	return s
}

func (s *Telegram) ScrapeStory(ctx context.Context, url string) (string, error) {
	return s.ScrapeStoryFrom(ctx, url, "div[id='article-content']", "p:not([class])", false)
}
