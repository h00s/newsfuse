package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/newsfuse/app/models"
	"github.com/uptrace/bun"
)

type TopicsService struct {
	raptor.Service
	Cache *CacheService
}

func (ts *TopicsService) All() (models.Topics, error) {
	var topics models.Topics

	if err := ts.memstoreGetTopics(&topics); err == nil {
		return topics, nil
	}

	err := ts.DB.Conn().(*bun.DB).
		NewSelect().
		Model(&topics).
		Order("id").
		Scan(context.Background())
	if err != nil {
		ts.Log.Error("Error geting topics", "Error", err.Error())
		return topics, errs.NewErrorInternal(err.Error())
	}

	go ts.memstoreSetTopics(&topics)
	return topics, nil
}

func (ts *TopicsService) memstoreGetTopics(topics *models.Topics) error {
	if data, ok := ts.Cache.Get("topics"); ok {
		json.Unmarshal(data, topics)
		return nil
	}
	ts.Log.Warn("Topics not found in memstore")
	return errors.New("topics not found in memstore")
}

func (ts *TopicsService) memstoreSetTopics(topics *models.Topics) {
	data, err := json.Marshal(*topics)
	if err != nil {
		ts.Log.Warn("Error setting topics in memstore", "error", err.Error())
	}
	ts.Cache.Set("topics", data)
}
