package scrapers

import (
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
)

type Mojportalhr struct {
	internal.DefaultScraper
	Source models.Source
}

func NewMojportalhr(h chan (models.Headlines), sourceID uint) *Mojportalhr {
	s := &Mojportalhr{
		DefaultScraper: *internal.NewScraper("MojPortal.hr", "https://www.mojportal.hr/", 15, 30, h),
	}

	s.ScrapeHeadline("div[class^='article_teaser__horizontal_box2']", func(e *colly.HTMLElement) {
		s.AddHeadline(models.Headline{
			SourceID:    sourceID,
			Title:       strings.TrimSpace(e.ChildText("span[class*='article_teaser__title_text']")),
			URL:         "https://www.mojportal.hr" + strings.TrimSpace(e.ChildAttr("a[class*='article_teaser__title_link']", "href")),
			PublishedAt: time.Now(),
		})
	})

	return s
}

func (s *Mojportalhr) ScrapeStory(url string) (string, error) {
	return s.DefaultScraper.ScrapeStory(url, "div[class='article__container']", "p:not([class])")
}
