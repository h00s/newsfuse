package robot

import (
	"github.com/h00s/newsfuse/api/mmc/models"
	"github.com/h00s/newsfuse/api/mmc/models/sources"
)

type Robot struct {
	Scrapers []models.Scraper
}

func NewRobot() *Robot {
	return &Robot{
		Scrapers: []models.Scraper{
			sources.NewRadioDaruvar(),
		},
	}
}

func (r *Robot) Start() {
	// TODO: Start scrapers
}
