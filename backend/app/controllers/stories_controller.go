package controllers

import (
	"strconv"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/newsfuse/app/services"
)

type StoriesController struct {
	raptor.Controller

	Stories *services.StoriesService
}

func (c *StoriesController) Get(ctx *raptor.Context) error {
	headlineID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return errs.NewErrorBadRequest("Invalid Headline ID")
	}

	story, err := c.Stories.Get(headlineID)
	if err != nil {
		return errs.NewErrorNotFound("Story not found")
	}

	return ctx.Data(story)
}

func (c *StoriesController) Summarize(ctx *raptor.Context) error {
	storyID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return errs.NewErrorBadRequest("Invalid Story ID")
	}

	story, err := c.Stories.Summarize(storyID)
	if err != nil {
		return errs.NewErrorNotFound("Story not found")
	}

	return ctx.Data(story)
}
