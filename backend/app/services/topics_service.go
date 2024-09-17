package services

import (
	"context"
	"encoding/json"

	"github.com/go-raptor/raptor/v3"
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

	err = ts.DB.
		NewSelect().
		Model(&topics).
		Order("id").
		Scan(context.Background())
	if err != nil {
		ts.Log.Error(err.Error())
	}

	go ts.Memstore.Set("topics", topics)
	return topics
}
