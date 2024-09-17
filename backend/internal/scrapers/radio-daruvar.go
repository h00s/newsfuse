package scrapers

import (
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
)

type RadioDaruvar struct {
	internal.DefaultScraper
}

func NewRadioDaruvar(h chan (models.Headlines), sourceID int64) *RadioDaruvar {
	s := &RadioDaruvar{
		DefaultScraper: *internal.NewScraper(h,
			"Radio Daruvar",
			"https://www.radio-daruvar.hr/",
			30,
			60,
			[]int{0, 1, 2, 3, 4, 5, 19, 20, 21, 22, 23},
		),
	}

	s.ScrapeHeadline("li[class='news-item'],h2[class='post-title']", func(e *colly.HTMLElement) {
		s.AddHeadline(models.Headline{
			SourceID:    sourceID,
			Title:       strings.TrimSpace(e.ChildText("a")),
			URL:         strings.TrimSpace(e.ChildAttr("a", "href")),
			PublishedAt: time.Now(),
		})
	})

	return s
}

func (s *RadioDaruvar) ScrapeStory(url string) (string, error) {
	return s.DefaultScraper.ScrapeStory(url, "div[class='entry-content entry clearfix']", "p:not([class])", false)
}
