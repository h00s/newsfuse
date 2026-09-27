package scrapers

import (
	"context"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/app/utils"
)

type N1InfoCroatia struct {
	utils.DefaultScraper
}

func NewN1InfoCroatia(sourceID int64) *N1InfoCroatia {
	s := &N1InfoCroatia{
		DefaultScraper: *utils.NewScraper(
			"N1",
			"https://n1info.hr/vijesti/",
			10,
			20,
			[]int{0, 1, 2, 3, 4, 5, 23},
		),
	}

	s.ScrapeHeadline("h3[data-testid='article-title']", func(e *colly.HTMLElement) {
		s.AddHeadline(models.Headline{
			SourceID:    sourceID,
			Title:       e.ChildText("a"),
			URL:         "https://n1info.hr" + e.ChildAttr("a", "href"),
			PublishedAt: time.Now(),
		})
	})

	return s
}

func (s *N1InfoCroatia) ScrapeStory(ctx context.Context, url string) (string, error) {
	return s.ScrapeStoryFrom(ctx, url, "div.article-content-wrapper", "p[data-block-key]", false)
}
