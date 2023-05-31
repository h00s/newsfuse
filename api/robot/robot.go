package robot

import "github.com/h00s/newsfuse/api/robot/sources"

type Robot struct {
	Scrapers []Scraper
}

func NewRobot() *Robot {
	return &Robot{
		Scrapers: []Scraper{
			sources.NewRadioDaruvar(),
		},
	}
}

func (r *Robot) Start() {
	// TODO: Start scrapers
}
