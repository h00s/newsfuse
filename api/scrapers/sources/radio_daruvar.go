package scrapers

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/api/mmc/models"
)

type RadioDaruvar struct {
	ID      int    `json:"id"`
	BaseURL string `json:"base_url"`
}

func NewRadioDaruvar() *RadioDaruvar {
	return &RadioDaruvar{
		ID:      1,
		BaseURL: "http://www.radio-daruvar.hr/",
	}
}

func (r *RadioDaruvar) GetHeadlines() ([]models.Headline, error) {
	var headlines []models.Headline
	c := colly.NewCollector()

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		headline := models.Headline{
			Title: e.Text,
			URL:   e.Attr("href"),
		}
		headlines = append(headlines, headline)
		fmt.Println(headline)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(r.BaseURL)
	return nil, nil
}
