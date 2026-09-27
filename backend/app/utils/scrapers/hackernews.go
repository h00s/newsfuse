package scrapers

import (
	"context"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/app/utils"
)

type HackerNews struct {
	utils.DefaultScraper
}

func NewHackerNews(sourceID int64) *HackerNews {
	s := &HackerNews{
		DefaultScraper: *utils.NewScraper(
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

func (s *HackerNews) ScrapeStory(ctx context.Context, url string) (string, error) {
	return s.ScrapeStoryFrom(ctx, url, "td[class='title']", "span[class='titleline']", true)
}
