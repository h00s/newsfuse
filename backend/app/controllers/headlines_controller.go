package controllers

import (
	"time"

	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/services"
)

type HeadlinesController struct {
	raptor.Controller
	Headlines *services.HeadlinesService
}

func (hc *HeadlinesController) All(c *raptor.Context) error {
	topicID, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(400)
	}
	if lastID := c.QueryInt("last_id", 0); lastID != 0 {
		return c.JSON(hc.Headlines.AllByLastID(topicID, lastID))
	}
	return c.JSON(hc.Headlines.All(topicID))
}

func (hc *HeadlinesController) Search(c *raptor.Context) error {
	query := c.Query("query")
	if query == "" || len(query) < 3 {
		return c.SendStatus(400)
	}
	return c.JSON(hc.Headlines.Search(query))
}

func (hc *HeadlinesController) Count(c *raptor.Context) error {
	topicID, _ := c.ParamsInt("id")
	status := c.Query("status")
	since := c.QueryInt("since", 0)
	if status != "" && since != 0 {
		sinceTime := time.Unix(int64(since/1000), 0)
		return c.JSON(
			raptor.Map{
				"count": hc.Headlines.Count(topicID, sinceTime),
			},
		)
	}

	return c.JSON("", 400)
}
