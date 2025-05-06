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

func (hc *HeadlinesController) All(c *raptor.Context) error {
	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return errs.NewErrorBadRequest("Invalid Topic ID")
	}

	var headlines models.Headlines
	if lastID, err := strconv.ParseInt(c.QueryParam("last_id"), 10, 64); err == nil {
		headlines, err = hc.Headlines.AllByLastID(topicID, lastID)
		if err != nil {
			return err
		}
		return c.Data(headlines)
	}

	headlines, err = hc.Headlines.All(topicID)
	if err != nil {
		return err
	}
	return c.Data(headlines)
}

func (hc *HeadlinesController) Search(c *raptor.Context) error {
	query := c.QueryParam("query")
	if query == "" || len(query) < 3 {
		return errs.NewErrorBadRequest("Invalid query")
	}

	var headlines models.Headlines
	headlines, err := hc.Headlines.Search(query)
	if err != nil {
		return err
	}
	return c.Data(headlines)
}

func (hc *HeadlinesController) Count(c *raptor.Context) error {
	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return errs.NewErrorBadRequest("Invalid Topic ID")
	}

	status := c.QueryParam("status")
	since, err := strconv.Atoi(c.QueryParam("since"))
	if err == nil && status != "" && since != 0 {
		sinceTime := time.Unix(int64(since/1000), 0)
		count, err := hc.Headlines.Count(topicID, sinceTime)
		if err != nil {
			return err
		}
		return c.Data(
			map[string]interface{}{
				"count": count,
			},
		)
	}

	return errs.NewErrorBadRequest("Invalid query parameters")
}
