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
	topicID, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(400)
	}
	return c.JSON(hc.Hs.All(topicID))
}

func (hc *HeadlinesController) Story(c *raptor.Context) error {
	headlineID, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(400)
	}

	story, err := hc.Hs.GetStory(headlineID)
	if err != nil {
		return c.SendStatus(404)
	}

	return c.JSON(story)
}
