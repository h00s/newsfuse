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

	s.Collector.OnHTML("li[class='news-item'],h2[class='post-title']", func(e *colly.HTMLElement) {
		s.AddHeadline(models.Headline{
			Source:      s.Name,
			Title:       strings.TrimSpace(e.ChildText("a")),
			URL:         strings.TrimSpace(e.ChildAttr("a", "href")),
			PublishedAt: time.Now(),
		})
	})

	s.Start()
	return s
}
