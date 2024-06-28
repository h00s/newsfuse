package scrapers

import (
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
)

type IndexhrWorld struct {
	internal.DefaultScraper
}

func NewIndexhrWorld(h chan (models.Headlines), sourceID uint) *IndexhrWorld {
	s := &IndexhrWorld{
		DefaultScraper: *internal.NewScraper(h,
			"Index.hr",
			"https://www.index.hr/vijesti/rubrika/hrvatska/23.aspx",
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

func (s *IndexhrWorld) ScrapeStory(url string) (string, error) {
	return s.DefaultScraper.ScrapeStory(url, "div[class^='text vijesti-link-underline']", "p:not([class])", false)
}
