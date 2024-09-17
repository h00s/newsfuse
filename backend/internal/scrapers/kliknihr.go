package scrapers

import (
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
)

type Kliknihr struct {
	internal.DefaultScraper
}

func NewKliknihr(h chan (models.Headlines), sourceID int64) *Kliknihr {
	s := &Kliknihr{
		DefaultScraper: *internal.NewScraper(h,
			"klikni.hr",
			"https://www.klikni.hr",
			15,
			30,
			[]int{0, 1, 2, 3, 4, 5, 23},
		),
	}

	s.ScrapeHeadline("h3.jeg_post_title", func(e *colly.HTMLElement) {
		s.AddHeadline(models.Headline{
			SourceID:    sourceID,
			Title:       e.ChildText("a"),
			URL:         e.ChildAttr("a", "href"),
			PublishedAt: time.Now(),
		})
	})

	return s
}

func (s *Kliknihr) ScrapeStory(url string) (string, error) {
	return s.DefaultScraper.ScrapeStory(url, "section[class='container page-content']", "p:not([class])", false)
}
