package controllers

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/app/services"
)

type TopicsController struct {
	raptor.Controller

	Topics *services.TopicsService
}

func (c *TopicsController) All(ctx *raptor.Context) error {
	topics, err := c.Topics.All()
	if err != nil {
		return err
	}

	return ctx.Data(topics)
}
