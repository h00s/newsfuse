package scrapers

import (
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
)

type HackerNews struct {
	internal.DefaultScraper
}

func NewHackerNews(h chan (models.Headlines), sourceID int64) *HackerNews {
	s := &HackerNews{
		DefaultScraper: *internal.NewScraper(h,
			"Hacker News",
			"https://news.ycombinator.com/",
			10,
			15,
			[]int{},
		),
	}

	s.ScrapeHeadline("tr[class^='athing']", func(e *colly.HTMLElement) {
		url := s.URL + "item?id=" + e.Attr("id")
		anchor := e.DOM.Find("span[class='titleline'] > a").First()
		s.AddHeadline(models.Headline{
			SourceID:    sourceID,
			Title:       anchor.Text(),
			URL:         url,
			PublishedAt: time.Now(),
		})
	})

	return s
}

func (s *HackerNews) ScrapeStory(url string) (string, error) {
	return s.DefaultScraper.ScrapeStory(url, "td[class='title']", "span[class='titleline']", true)
}
