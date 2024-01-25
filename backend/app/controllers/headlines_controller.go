package controllers

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/services"
)

type HeadlinesController struct {
	raptor.Controller

	Hs *services.HeadlinesService
}

func (hc *HeadlinesController) Index(c *raptor.Context) error {
	return c.JSON(hc.Hs.Headlines)
}
