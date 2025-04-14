package controllers

import (
	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/newsfuse/app/services"
)

type TopicsController struct {
	raptor.Controller
	Topics *services.TopicsService
}

func (sc *TopicsController) All(c raptor.State) error {
	topics, err := sc.Topics.All()
	if err != nil {
		return err
	}
	return c.JSONResponse(topics)
}
