package controllers

import (
	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/newsfuse/app/services"
)

type SourcesController struct {
	raptor.Controller
	Sources *services.SourcesService
}

func (sc *SourcesController) All(c *raptor.Context) error {
	sources, err := sc.Sources.All()
	if err != nil {
		return c.JSONError(err)
	}
	return c.JSONResponse(sources)
}
