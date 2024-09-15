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

func (ss *SourcesService) All() models.Sources {
	var sources models.Sources
	data, err := ss.Memstore.Get("sources")
	if err == nil && data != "" {
		json.Unmarshal([]byte(data), &sources)
		return sources
	}
	ss.DB.
		NewSelect().
		Model(&sources).
		Exec(context.Background())
	ss.Memstore.Set("sources", sources)
	return sources
}

func (ss *SourcesService) Get(id uint) models.Source {
	var source models.Source

	ss.DB.
		NewSelect().
		Model(&source).
		Where("id = ?", id).
		Exec(context.Background())

	return source
}
