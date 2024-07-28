package scrapers

import (
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
)

type Telegram struct {
	internal.DefaultScraper
}

func NewTelegram(h chan (models.Headlines), sourceID uint) *Telegram {
	s := &Telegram{
		DefaultScraper: *internal.NewScraper(h,
			"Telegram",
			"https://www.telegram.hr/vijesti",
			15,
			30,
			[]int{1, 2, 3, 4, 5},
		),
	}

	s.ScrapeHeadline("a[role='article']", func(e *colly.HTMLElement) {
		s.AddHeadline(models.Headline{
			SourceID:    sourceID,
			Title:       e.ChildText("h2"),
			URL:         "https://www.telegram.hr" + e.Attr("href"),
			PublishedAt: time.Now(),
		})
	})

	return s
}

func (s *Telegram) ScrapeStory(url string) (string, error) {
	return s.DefaultScraper.ScrapeStory(url, "div[id='article-content']", "p:not([class])", false)
}
