package controllers

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/services"
)

type HeadlinesController struct {
	raptor.Controller
}

func (hc *HeadlinesController) Headlines() *services.HeadlinesService {
	return hc.Services["HeadlinesService"].(*services.HeadlinesService)
}

func (hc *HeadlinesController) All(c *raptor.Context) error {
	topicID, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(400)
	}
	return c.JSON(hc.Headlines().All(topicID))
}

func (hc *HeadlinesController) Story(c *raptor.Context) error {
	headlineID, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(400)
	}

	story, err := hc.Headlines().GetStory(headlineID)
	if err != nil {
		return c.SendStatus(404)
	}

	return c.JSON(story)
}
