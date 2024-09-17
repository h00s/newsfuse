package services

import (
	"context"
	"encoding/json"

	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/newsfuse/app/models"
)

type SourcesService struct {
	raptor.Service
	Memstore *Memstore
}

func (ss *SourcesService) All() (models.Sources, error) {
	var sources models.Sources
	data, err := ss.Memstore.Get("sources")
	if err == nil && data != "" {
		json.Unmarshal([]byte(data), &sources)
		return sources, nil
	}

	err = ss.DB.
		NewSelect().
		Model(&sources).
		Scan(context.Background())
	if err != nil {
		ss.Log.Error(err.Error())
		return sources, raptor.NewErrorInternal(err.Error())
	}
	go ss.Memstore.Set("sources", sources)
	return sources, nil
}

func (ss *SourcesService) Get(id int64) models.Source {
	var source models.Source

	ss.DB.
		NewSelect().
		Model(&source).
		Where("id = ?", id).
		Scan(context.Background())

	return source
}
