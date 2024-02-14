package controllers

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/services"
)

type StoriesController struct {
	raptor.Controller
}

func (sc *StoriesController) Stories() *services.StoriesService {
	return sc.Services["StoriesService"].(*services.StoriesService)
}

func (sc *StoriesController) Get(c *raptor.Context) error {
	headlineID, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(400)
	}

	story, err := sc.Stories().Get(headlineID)
	if err != nil {
		return c.SendStatus(404)
	}

	return c.JSON(story)
}
