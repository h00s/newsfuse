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

func (s *TopicsService) All() (models.Topics, error) {
	var topics models.Topics

	if err := s.memstoreGetTopics(&topics); err == nil {
		return topics, nil
	}

	err := s.Database.Conn().(*bun.DB).
		NewSelect().
		Model(&topics).
		Order("id").
		Scan(context.Background())
	if err != nil {
		s.Log.Error("Error geting topics", "Error", err.Error())
		return topics, errs.NewErrorInternal(err.Error())
	}

	go s.memstoreSetTopics(&topics)
	return topics, nil
}

func (s *TopicsService) memstoreGetTopics(topics *models.Topics) error {
	if data, ok := s.Cache.Get("topics"); ok {
		json.Unmarshal(data, topics)
		return nil
	}
	s.Log.Warn("Topics not found in memstore")
	return errors.New("topics not found in memstore")
}

func (s *TopicsService) memstoreSetTopics(topics *models.Topics) {
	data, err := json.Marshal(*topics)
	if err != nil {
		s.Log.Warn("Error setting topics in memstore", "error", err.Error())
	}
	s.Cache.Set("topics", data)
}
