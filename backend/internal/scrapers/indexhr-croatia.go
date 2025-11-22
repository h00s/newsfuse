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

func NewIndexhrCroatia(h chan (models.Headlines), sourceID int64) *IndexhrCroatia {
	s := &IndexhrCroatia{
		DefaultScraper: *internal.NewScraper(h,
			"Index.hr",
			"https://www.index.hr/vijesti/rubrika/hrvatska/22.aspx",
			5,
			15,
			[]int{1, 2, 3, 4, 5},
		),
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
	return s.DefaultScraper.ScrapeStory(url, "div[class='text-holder']", "p:not([class])", false)
}
