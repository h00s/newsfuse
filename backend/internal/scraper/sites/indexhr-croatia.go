package sites

import (
	"context"
	"strings"

	"github.com/h00s/newsfuse/internal/collector"
	"github.com/h00s/newsfuse/internal/scraper"
)

type IndexhrCroatia struct {
	scraper.Base
}

func NewIndexhrCroatia() *IndexhrCroatia {
	s := &IndexhrCroatia{
		Base: *scraper.New(
			"Index.hr",
			"https://www.index.hr/vijesti/rubrika/hrvatska/22.aspx",
			5,
			15,
			[]int{1, 2, 3, 4, 5},
		),
	}

	s.ScrapeHeadline("a[class='vijesti-text-hover scale-img-hover']", func(e *collector.HTMLElement) {
		url := e.Attr("href")
		if index := strings.Index(url, "?"); index != -1 {
			url = url[:index]
		}

		s.AddHeadline(scraper.Headline{
			Title: e.ChildText(".title"),
			URL:   "https://www.index.hr" + url,
		})
	})

	return s
}

func (s *IndexhrCroatia) ScrapeStory(ctx context.Context, url string) (string, error) {
	return s.ScrapeStoryFrom(ctx, url, "div[class='text-holder']", "p:not([class])", false)
}
