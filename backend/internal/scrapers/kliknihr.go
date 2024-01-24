package scrapers

import (
	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/internal"
)

type Kliknihr struct {
	internal.DefaultScraper
}

func NewKliknihr() internal.Scraper {
	s := internal.NewScraper("kliknihr", "https://www.kliknihr.net/", 15, 75)
	s.Collector.OnHTML("div#content div.post", func(e *colly.HTMLElement) {
	})
	return s
}
