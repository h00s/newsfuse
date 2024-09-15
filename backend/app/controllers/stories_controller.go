package controllers

import (
	"strconv"

	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/newsfuse/app/services"
)

type StoriesController struct {
	raptor.Controller
	Stories *services.StoriesService
}

func (sc *StoriesController) Get(c *raptor.Context) error {
	headlineID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSONError(raptor.NewErrorBadRequest("Invalid Headline ID"))
	}

	story, err := sc.Stories.Get(headlineID)
	if err != nil {
		return c.JSONError(raptor.NewErrorNotFound("Story not found"))
	}

	return c.JSON(story)
}

func (sc *StoriesController) Summarize(c *raptor.Context) error {
	//storyID, err := c.ParamsInt("id")
	storyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSONError(raptor.NewErrorBadRequest("Invalid Story ID"))
	}

	story, err := sc.Stories.Summarize(storyID)
	if err != nil {
		return c.JSONError(raptor.NewErrorNotFound("Story not found"))
	}

	return c.JSON(story)
}
