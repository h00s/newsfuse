package controllers

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/app/services"
)

type SourcesController struct {
	raptor.Controller

	Sources *services.SourcesService
}

// Index lists every source, or one topic's when ?topicId= is given.
func (c *SourcesController) Index(ctx *raptor.Context) error {
	topicID, filtered, err := queryID(ctx, "topicId")
	if err != nil {
		return err
	}

	var sources models.Sources
	if filtered {
		sources, err = c.Sources.ListByTopic(topicID)
	} else {
		sources, err = c.Sources.List()
	}
	if err != nil {
		return err
	}
	return ctx.Data(models.NewSourceResponses(sources))
}
