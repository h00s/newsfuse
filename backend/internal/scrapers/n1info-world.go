package scrapers

import (
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
)

type N1InfoWorld struct {
	internal.DefaultScraper
}

func NewN1InfoWorld(h chan (models.Headlines), sourceID int64) *N1InfoWorld {
	s := &N1InfoWorld{
		DefaultScraper: *internal.NewScraper(h,
			"N1",
			"https://n1info.hr/svijet/",
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

func (s *N1InfoWorld) ScrapeStory(url string) (string, error) {
	return s.DefaultScraper.ScrapeStory(url, "div.article-content-wrapper", "p[data-block-key]", false)
}
