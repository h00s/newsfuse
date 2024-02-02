package controllers

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/services"
)

type SourcesController struct {
	raptor.Controller

	Ss *services.SourcesService
}

func (sc *SourcesController) All(c *raptor.Context) error {
	return c.JSON(sc.Ss.All())
}
