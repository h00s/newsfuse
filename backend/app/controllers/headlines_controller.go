package controllers

import (
	"strconv"
	"time"

	"github.com/go-raptor/components"
	"github.com/go-raptor/errs"
	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/app/services"
)

type HeadlinesController struct {
	raptor.Controller
	Headlines *services.HeadlinesService
}

func (hc *HeadlinesController) All(c *components.Context) error {
	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSONError(errs.NewErrorBadRequest("Invalid Topic ID"))
	}

	var headlines models.Headlines
	if lastID, err := strconv.ParseInt(c.QueryParam("last_id"), 10, 64); err == nil {
		headlines, err = hc.Headlines.AllByLastID(topicID, lastID)
		if err == nil {
			return c.JSON(headlines)
		}
		return c.JSONError(err)
	}

	headlines, err = hc.Headlines.All(topicID)
	if err == nil {
		return c.JSON(headlines)
	}
	return c.JSONError(err)
}

func (hc *HeadlinesController) Search(c *components.Context) error {
	query := c.QueryParam("query")
	if query == "" || len(query) < 3 {
		return c.JSONError(errs.NewErrorBadRequest("Invalid query"))
	}

	var headlines models.Headlines
	headlines, err := hc.Headlines.Search(query)
	if err == nil {
		return c.JSON(headlines)
	}
	return c.JSONError(err)
}

func (hc *HeadlinesController) Count(c *components.Context) error {
	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSONError(errs.NewErrorBadRequest("Invalid Topic ID"))
	}
	status := c.QueryParam("status")
	since, err := strconv.Atoi(c.QueryParam("since"))
	if err == nil && status != "" && since != 0 {
		sinceTime := time.Unix(int64(since/1000), 0)
		count, err := hc.Headlines.Count(topicID, sinceTime)
		if err == nil {
			return c.JSON(
				raptor.Map{
					"count": count,
				},
			)
		}
		return c.JSONError(err)
	}

	return c.JSONError(errs.NewErrorBadRequest("Invalid query parameters"))
}
