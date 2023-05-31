package robot

import (
	"github.com/h00s/newsfuse/api/mmc/models"
)

type Scraper interface {
	GetHeadlines() ([]models.Headline, error)
	GetInterval() int
}
