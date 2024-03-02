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

func NewHackerNews(h chan (models.Headline), sourceID uint) *HackerNews {
	s := &HackerNews{
		DefaultScraper: *internal.NewScraper("Hacker News", "https://news.ycombinator.com/", 10, 15, h),
	}

	s.ScrapeHeadline("span[class='titleline']", func(e *colly.HTMLElement) {
		anchor := e.DOM.Find("a").First()
		url, _ := anchor.Attr("href")
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
	return "Scraping story from Hacker News is not available", nil
}
