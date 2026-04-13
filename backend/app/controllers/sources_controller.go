package controllers

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/app/services"
)

type SourcesController struct {
	raptor.Controller

	Sources *services.SourcesService
}

func (c *SourcesController) All(ctx *raptor.Context) error {
	sources, err := c.Sources.All()
	if err != nil {
		return err
	}

	return ctx.Data(sources)
}
