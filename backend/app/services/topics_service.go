package services

import (
	"encoding/json"

	"github.com/go-raptor/raptor/v2"
	"github.com/h00s/newsfuse/app/models"
)

type TopicsService struct {
	raptor.Service
	Memstore *Memstore
}

func (ts *TopicsService) All() models.Topics {
	var topics models.Topics
	data, err := ts.Memstore.Get("topics")
	if err == nil && data != "" {
		json.Unmarshal([]byte(data), &topics)
		return topics
	}
	ts.DB.Order("id").Find(&topics)
	ts.Memstore.Set("topics", topics)
	return topics
}
