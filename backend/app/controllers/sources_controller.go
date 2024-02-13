package controllers

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/services"
)

type SourcesController struct {
	raptor.Controller
}

func (sc *SourcesController) Sources() *services.SourcesService {
	return sc.Services["SourcesService"].(*services.SourcesService)
}

func (sc *SourcesController) All(c *raptor.Context) error {
	return c.JSON(sc.Sources().All())
}
