package services

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/models"
)

type TopicsService struct {
	raptor.Service
}

func (ts *TopicsService) All() models.Topics {
	var topics models.Topics
	ts.DB.Order("id").Find(&topics)
	return topics
}
