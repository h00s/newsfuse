// Package scrapers implmements all scrapers for different news sources.
package scrapers

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/app/utils"
)

type Bughr struct {
	utils.DefaultScraper
}

func NewBughr(h chan models.Headlines, log *slog.Logger, sourceID int64) *Bughr {
	s := &Bughr{
		DefaultScraper: *utils.NewScraper(h, log,
			"Bug",
			"https://www.bug.hr",
			10,
			30,
			[]int{0, 1, 2, 3, 4, 5, 23},
		),
	}

	s.ScrapeHeadline("h2[class='post-listing__title'] > a", func(e *colly.HTMLElement) {
		// skip premium articles
		if e.DOM.Find("svg").Length() > 0 {
			return
		}

		url := e.Attr("href")
		// skip forum articles
		if strings.HasPrefix(url, "/forum") {
			return
		}
		if !strings.HasPrefix(url, "http") {
			url = s.URL + url
		}

		s.AddHeadline(models.Headline{
			SourceID:    sourceID,
			Title:       e.Text,
			URL:         url,
			PublishedAt: time.Now(),
		})
	})

	return s
}

func (s *Bughr) ScrapeStory(url string) (string, error) {
	return s.DefaultScraper.ScrapeStory(url, "div[class^='post-full__content']", "p:not([class])", false)
}
