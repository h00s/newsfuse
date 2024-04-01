package controllers

import (
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
