package controllers

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/services"
)

type HeadlinesController struct {
	raptor.Controller

	Hs *services.HeadlinesService
}

func (hc *HeadlinesController) All(c *raptor.Context) error {
	return c.JSON(hc.Hs.All())
}

func (hc *HeadlinesController) Story(c *raptor.Context) error {
	headlineID, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(400)
	}

	story, err := hc.Hs.Story(headlineID)
	if err != nil {
		return c.SendStatus(404)
	}

	return c.JSON(story)
}
