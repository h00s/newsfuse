package scrapers

import (
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
)

type Bughr struct {
	internal.DefaultScraper
}

func NewBughr(h chan (models.Headlines), sourceID uint) *Bughr {
	s := &Bughr{
		DefaultScraper: *internal.NewScraper("Bug", "https://www.bug.hr", 10, 30, h),
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
