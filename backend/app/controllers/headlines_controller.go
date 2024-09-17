package controllers

import (
	"strconv"
	"time"

	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/newsfuse/app/services"
)

type HeadlinesController struct {
	raptor.Controller
	Headlines *services.HeadlinesService
}

func (hc *HeadlinesController) All(c *raptor.Context) error {
	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSONError(raptor.NewErrorBadRequest("Invalid Topic ID"))
	}
	if lastID, err := strconv.ParseInt(c.QueryParam("last_id"), 10, 64); err == nil {
		return c.JSON(hc.Headlines.AllByLastID(topicID, lastID))
	}
	return c.JSON(hc.Headlines.All(topicID))
}

func (hc *HeadlinesController) Search(c *raptor.Context) error {
	query := c.QueryParam("query")
	if query == "" || len(query) < 3 {
		return c.JSONError(raptor.NewErrorBadRequest("Invalid query"))
	}
	return c.JSON(hc.Headlines.Search(query))
}

func (hc *HeadlinesController) Count(c *raptor.Context) error {
	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSONError(raptor.NewErrorBadRequest("Invalid Topic ID"))
	}
	status := c.QueryParam("status")
	since, err := strconv.Atoi(c.QueryParam("since"))
	if err == nil && status != "" && since != 0 {
		sinceTime := time.Unix(int64(since/1000), 0)
		return c.JSON(
			raptor.Map{
				"count": hc.Headlines.Count(topicID, sinceTime),
			},
		)
	}

	return c.JSONError(raptor.NewErrorBadRequest("Invalid query parameters"))
}
