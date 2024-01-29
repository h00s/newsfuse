package scrapers

import (
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
)

type RadioDaruvar struct {
	internal.DefaultScraper
}

func NewRadioDaruvar(headline chan (models.Headline)) internal.Scraper {
	s := internal.NewScraper(headline, "Radio Daruvar", "https://www.radio-daruvar.hr/", 30, 60)

	s.Collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
		s.Headlines = models.Headlines{}
	})

	s.Collector.OnHTML("li[class='news-item'],h2[class='post-title']", func(e *colly.HTMLElement) {
		headline := models.Headline{
			Source:      s.Name,
			Title:       strings.TrimSpace(e.ChildText("a")),
			URL:         strings.TrimSpace(e.ChildAttr("a", "href")),
			PublishedAt: time.Now(),
		}
		s.Headlines = append(s.Headlines, headline)
	})

	s.Collector.OnScraped(func(r *colly.Response) {
		for i := len(s.Headlines) - 1; i >= 0; i-- {
			h := s.Headlines[i]
			s.HeadlineChannel <- h
		}
	})

	s.Start()
	return s
}
