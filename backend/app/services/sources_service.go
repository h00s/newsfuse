package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/newsfuse/app/models"
	"github.com/uptrace/bun"
)

type SourcesService struct {
	raptor.Service

	Cache *CacheService
}

func (s *SourcesService) All() (models.Sources, error) {
	var sources models.Sources

	if err := s.memstoreGetSources(&sources); err == nil {
		return sources, nil
	}

	err := s.Database.Conn().(*bun.DB).
		NewSelect().
		Model(&sources).
		Scan(context.Background())
	if err != nil {
		s.Log.Error(err.Error())
		return sources, errs.NewErrorInternal(err.Error())
	}

	go s.memstoreSetSources(&sources)
	return sources, nil
}

func (s *SourcesService) Get(id int64) models.Source {
	var source models.Source

	if err := s.memstoreGetSource(id, &source); err == nil {
		return source
	}

	s.Database.Conn().(*bun.DB).
		NewSelect().
		Model(&source).
		Where("id = ?", id).
		Scan(context.Background())

	go s.memstoreSetSource(&source)
	return source
}

func (s *SourcesService) memstoreGetSources(sources *models.Sources) error {
	if data, ok := s.Cache.Get("sources"); ok {
		json.Unmarshal(data, sources)
		return nil
	}

	s.Log.Warn("Sources not found in memstore")
	return errors.New("sources not found in memstore")
}

func (s *SourcesService) memstoreSetSources(sources *models.Sources) {
	data, err := json.Marshal(*sources)
	if err != nil {
		s.Log.Warn("Error setting sources in memstore", "error", err.Error())
	}
	s.Cache.Set("sources", data)
}

func (s *SourcesService) memstoreGetSource(id int64, source *models.Source) error {
	if data, ok := s.Cache.Get(fmt.Sprintf("sources:%d", id)); ok {
		json.Unmarshal(data, source)
		return nil
	}

	s.Log.Warn("Source not found in memstore", "source", id)
	return errors.New("source not found in memstore")
}

func (s *SourcesService) memstoreSetSource(source *models.Source) {
	data, err := json.Marshal(*source)
	if err != nil {
		s.Log.Warn("Error setting source in memstore", "error", err.Error())
	}
	s.Cache.Set(fmt.Sprintf("sources:%d", source.ID), data)
}
