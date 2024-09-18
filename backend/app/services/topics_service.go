package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/newsfuse/app/models"
)

type TopicsService struct {
	raptor.Service
	Memstore *Memstore
}

func (ts *TopicsService) All() (models.Topics, error) {
	var topics models.Topics

	if err := ts.memstoreGetTopics(&topics); err == nil {
		return topics, nil
	}

	err := ts.DB.
		NewSelect().
		Model(&topics).
		Order("id").
		Scan(context.Background())
	if err != nil {
		ts.Log.Error("Error geting topics", "Error", err.Error())
		return topics, raptor.NewErrorInternal(err.Error())
	}

	go ts.memstoreSetTopics(&topics)
	return topics, nil
}

func (ts *TopicsService) memstoreGetTopics(topics *models.Topics) error {
	data, err := ts.Memstore.Get("topics")
	if err == nil && data != "" {
		json.Unmarshal([]byte(data), topics)
		return nil
	}

	ts.Log.Warn("Topics not found in memstore")
	return errors.New("topics not found in memstore")
}

func (ts *TopicsService) memstoreSetTopics(topics *models.Topics) {
	err := ts.Memstore.Set("topics", *topics)
	if err != nil {
		ts.Log.Warn("Error setting topics in memstore", "error", err.Error())
	}
}
