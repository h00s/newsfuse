package controllers

import (
	"github.com/go-raptor/raptor/v2"
	"github.com/h00s/newsfuse/app/services"
)

type TopicsController struct {
	raptor.Controller
	Topics *services.TopicsService
}

func (sc *TopicsController) All(c *raptor.Context) error {
	return c.JSON(sc.Topics.All())
}
