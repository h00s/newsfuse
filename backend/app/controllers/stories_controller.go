package controllers

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/app/services"
)

type StoriesController struct {
	raptor.Controller

	Stories *services.StoriesService
}

// Summarize summarizes a story with the LLM, once; later calls return the stored summary.
func (c *StoriesController) Summarize(ctx *raptor.Context) error {
	id, err := ctx.ParamInt64("id")
	if err != nil {
		return err
	}
	story, err := c.Stories.Summarize(ctx.Request().Context(), id)
	if err != nil {
		return err
	}
	return ctx.Data(models.NewStoryResponse(story))
}
