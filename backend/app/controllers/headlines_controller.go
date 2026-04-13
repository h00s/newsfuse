// Package controllers contains all newsfuse controllers
package controllers

import (
	"strconv"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/app/services"
)

type HeadlinesController struct {
	raptor.Controller

	Headlines *services.HeadlinesService
}

func (c *HeadlinesController) All(ctx *raptor.Context) error {
	topicID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return errs.NewErrorBadRequest("Invalid Topic ID")
	}

	var headlines models.Headlines
	if lastID, err := strconv.ParseInt(ctx.QueryParam("last_id"), 10, 64); err == nil {
		headlines, err = c.Headlines.AllByLastID(topicID, lastID)
		if err != nil {
			return err
		}
		return ctx.Data(headlines)
	}

	headlines, err = c.Headlines.All(topicID)
	if err != nil {
		return err
	}
	return ctx.Data(headlines)
}

func (c *HeadlinesController) Search(ctx *raptor.Context) error {
	query := ctx.QueryParam("query")
	if query == "" || len(query) < 3 {
		return errs.NewErrorBadRequest("Invalid query")
	}

	var headlines models.Headlines
	headlines, err := c.Headlines.Search(query)
	if err != nil {
		return err
	}

	return ctx.Data(headlines)
}

func (c *HeadlinesController) Count(ctx *raptor.Context) error {
	topicID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return errs.NewErrorBadRequest("Invalid Topic ID")
	}

	status := ctx.QueryParam("status")
	since, err := strconv.Atoi(ctx.QueryParam("since"))
	if err == nil && status != "" && since != 0 {
		sinceTime := time.Unix(int64(since/1000), 0)
		count, err := c.Headlines.Count(topicID, sinceTime)
		if err != nil {
			return err
		}
		return ctx.Data(
			map[string]interface{}{
				"count": count,
			},
		)
	}

	return errs.NewErrorBadRequest("Invalid query parameters")
}
