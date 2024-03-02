package scrapers

import (
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
)

type N1InfoCroatia struct {
	internal.DefaultScraper
}

func NewN1InfoCroatia(h chan (models.Headline), sourceID uint) *N1InfoCroatia {
	s := &N1InfoCroatia{
		DefaultScraper: *internal.NewScraper("N1", "https://n1info.hr/vijesti/", 10, 20, h),
	}

	s.ScrapeHeadline("a[class='uc-block-post-grid-title-link']", func(e *colly.HTMLElement) {
		s.AddHeadline(models.Headline{
			SourceID:    sourceID,
			Title:       e.Text,
			URL:         e.Attr("href"),
			PublishedAt: time.Now(),
		})
	})

	return s
}

func (s *N1InfoCroatia) ScrapeStory(url string) (string, error) {
	return s.DefaultScraper.ScrapeStory(url, "div[class='entry-content']", "p:not([class])")
}
