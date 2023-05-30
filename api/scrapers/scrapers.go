package scrapers

import (
	"github.com/h00s/newsfuse/api/mmc/models"
	scrapers "github.com/h00s/newsfuse/api/scrapers/sources"
)

type Scraper interface {
	GetHeadlines() ([]models.Headline, error)
}

func NewScrapers() []Scraper {
	return []Scraper{
		scrapers.NewRadioDaruvar(),
	}
}
