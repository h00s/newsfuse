package controllers

import (
	"strconv"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/core"
	"github.com/h00s/newsfuse/app/services"
)

type StoriesController struct {
	raptor.Controller

	Stories *services.StoriesService
}

func (sc *StoriesController) Get(c *raptor.Context) error {
	headlineID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return core.NewErrorBadRequest("Invalid Headline ID")
	}

	story, err := sc.Stories.Get(headlineID)
	if err != nil {
		return core.NewErrorNotFound("Story not found")
	}

	return c.Data(story)
}

func (sc *StoriesController) Summarize(c *raptor.Context) error {
	storyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return core.NewErrorBadRequest("Invalid Story ID")
	}

	story, err := sc.Stories.Summarize(storyID)
	if err != nil {
		return core.NewErrorNotFound("Story not found")
	}

	return c.Data(story)
}
