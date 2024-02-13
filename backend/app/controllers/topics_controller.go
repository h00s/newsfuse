package controllers

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/services"
)

type TopicsController struct {
	raptor.Controller
}

func (sc *TopicsController) Topics() *services.TopicsService {
	return sc.Services["TopicsService"].(*services.TopicsService)
}

func (sc *TopicsController) All(c *raptor.Context) error {
	return c.JSON(sc.Topics().All())
}
