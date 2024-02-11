package controllers

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/services"
)

type TopicsController struct {
	raptor.Controller

	Ts *services.TopicsService
}

func (sc *TopicsController) All(c *raptor.Context) error {
	return c.JSON(sc.Ts.All())
}
