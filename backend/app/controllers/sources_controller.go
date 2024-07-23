package controllers

import (
	"github.com/go-raptor/raptor/v2"
	"github.com/h00s/newsfuse/app/services"
)

type SourcesController struct {
	raptor.Controller
	Sources *services.SourcesService
}

func (sc *SourcesController) All(c *raptor.Context) error {
	return c.JSON(sc.Sources.All())
}
