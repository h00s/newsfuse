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
}

func NewMojportalhr(headline chan (models.Headline)) internal.Scraper {
	s := internal.NewScraper(headline, "MojPortal.hr", "https://www.mojportal.hr/", 15, 45)

	s.Collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
		s.Headlines = models.Headlines{}
	})

	s.Collector.OnHTML("div[class^='article_teaser__horizontal_box2']", func(e *colly.HTMLElement) {
		headline := models.Headline{
			Source:      s.Name,
			Title:       strings.TrimSpace(e.ChildText("span[class*='article_teaser__title_text']")),
			URL:         "https://www.mojportal.hr/" + strings.TrimSpace(e.ChildAttr("a[class*='article_teaser__title_link']", "href")),
			PublishedAt: time.Now(),
		}
		s.Headlines = append(s.Headlines, headline)
	})

	s.Collector.OnScraped(func(r *colly.Response) {
		for i := len(s.Headlines) - 1; i >= 0; i-- {
			h := s.Headlines[i]
			s.Headline <- h
		}
	})

	s.Start()
	return s
}
