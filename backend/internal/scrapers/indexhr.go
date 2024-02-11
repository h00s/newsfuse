package scrapers

import (
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
)

/* https://www.index.hr/najnovije?kategorija=3 */

type Indexhr struct {
	internal.DefaultScraper
}

func NewIndexhr(h chan (models.Headline), sourceID uint) *Indexhr {
	s := &Indexhr{
		DefaultScraper: *internal.NewScraper("Index.hr", "https://www.index.hr/najnovije?kategorija=3", 5, 15, h),
	}

	s.ScrapeHeadline("div[class='title-box']", func(e *colly.HTMLElement) {
		url := e.ChildAttr("a[class='vijesti-text-hover scale-img-hover flex']", "href")
		if index := strings.Index(url, "?"); index != -1 {
			url = url[:index]
		}

		s.AddHeadline(models.Headline{
			SourceID:    sourceID,
			Title:       e.ChildText("h3[class='title']"),
			URL:         "https://www.index.hr" + url,
			PublishedAt: time.Now(),
		})
	})

	return s
}

func (s *Indexhr) ScrapeStory(url string) (string, error) {
	return s.DefaultScraper.ScrapeStory(url, "section[class='container page-content']", "p:not([class])")
}
