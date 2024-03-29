package scrapers

import (
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
)

type IndexhrCroatia struct {
	internal.DefaultScraper
}

func NewIndexhrCroatia(h chan (models.Headlines), sourceID uint) *IndexhrCroatia {
	s := &IndexhrCroatia{
		DefaultScraper: *internal.NewScraper("Index.hr", "https://www.index.hr/vijesti/rubrika/hrvatska/22.aspx", 5, 15, h),
	}

	s.ScrapeHeadline("a[class='vijesti-text-hover scale-img-hover']", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		if index := strings.Index(url, "?"); index != -1 {
			url = url[:index]
		}

		s.AddHeadline(models.Headline{
			SourceID:    sourceID,
			Title:       e.ChildText(".title"),
			URL:         "https://www.index.hr" + url,
			PublishedAt: time.Now(),
		})
	})

	return s
}

func (s *IndexhrCroatia) ScrapeStory(url string) (string, error) {
	return s.DefaultScraper.ScrapeStory(url, "div[class='text vijesti-link-underline']", "p:not([class])")
}
