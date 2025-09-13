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

func (ss *SourcesService) All() (models.Sources, error) {
	var sources models.Sources

	if err := ss.memstoreGetSources(&sources); err == nil {
		return sources, nil
	}

	err := ss.Database.Conn().(*bun.DB).
		NewSelect().
		Model(&sources).
		Scan(context.Background())
	if err != nil {
		ss.Log.Error(err.Error())
		return sources, errs.NewErrorInternal(err.Error())
	}
	go ss.memstoreSetSources(&sources)
	return sources, nil
}

func (ss *SourcesService) Get(id int64) models.Source {
	var source models.Source

	if err := ss.memstoreGetSource(id, &source); err == nil {
		return source
	}

	ss.Database.Conn().(*bun.DB).
		NewSelect().
		Model(&source).
		Where("id = ?", id).
		Scan(context.Background())

	go ss.memstoreSetSource(&source)
	return source
}

func (ss *SourcesService) memstoreGetSources(sources *models.Sources) error {
	if data, ok := ss.Cache.Get("sources"); ok {
		json.Unmarshal(data, sources)
		return nil
	}
	ss.Log.Warn("Sources not found in memstore")
	return errors.New("sources not found in memstore")
}

func (ss *SourcesService) memstoreSetSources(sources *models.Sources) {
	data, err := json.Marshal(*sources)
	if err != nil {
		ss.Log.Warn("Error setting sources in memstore", "error", err.Error())
	}
	ss.Cache.Set("sources", data)
}

func (ss *SourcesService) memstoreGetSource(id int64, source *models.Source) error {
	if data, ok := ss.Cache.Get(fmt.Sprintf("sources:%d", id)); ok {
		json.Unmarshal(data, source)
		return nil
	}
	ss.Log.Warn("Source not found in memstore", "source", id)
	return errors.New("source not found in memstore")
}

func (ss *SourcesService) memstoreSetSource(source *models.Source) {
	data, err := json.Marshal(*source)
	if err != nil {
		ss.Log.Warn("Error setting source in memstore", "error", err.Error())
	}
	ss.Cache.Set(fmt.Sprintf("sources:%d", source.ID), data)
}
