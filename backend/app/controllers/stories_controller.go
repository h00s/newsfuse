package controllers

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/services"
)

type StoriesController struct {
	raptor.Controller
	Stories *services.StoriesService
}

func (sc *StoriesController) Get(c *raptor.Context) error {
	headlineID, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(400)
	}

	story, err := sc.Stories.Get(headlineID)
	if err != nil {
		return c.SendStatus(404)
	}

	return c.JSON(story)
}

func (sc *StoriesController) Summarize(c *raptor.Context) error {
	storyID, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(400)
	}

	story, err := sc.Stories.Summarize(storyID)
	if err != nil {
		return c.SendStatus(404)
	}

	return c.JSON(story)
}
