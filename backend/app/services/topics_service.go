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

func (ts *TopicsService) All() (models.Topics, error) {
	var topics models.Topics

	data, err := ts.Memstore.Get("topics")
	if err == nil && data != "" {
		json.Unmarshal([]byte(data), &topics)
		return topics, nil
	}

	err = ts.DB.
		NewSelect().
		Model(&topics).
		Order("id").
		Scan(context.Background())
	if err != nil {
		ts.Log.Error("Error geting topics", "Error", err.Error())
		return topics, raptor.NewErrorInternal(err.Error())
	}

	go ts.Memstore.Set("topics", topics)
	return topics, nil
}
