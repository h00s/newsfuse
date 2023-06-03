package robot

import "github.com/h00s/newsfuse/api/mmc/middleware"

type Robot struct {
	models *middleware.ModelsMiddleware
}

func NewRobot(mm *middleware.ModelsMiddleware) *Robot {
	return &Robot{
		models: mm,
	}
}

func (r *Robot) Start() {
	// TODO: Start scrapers
}
