package services

import (
	"encoding/json"

	"github.com/go-raptor/raptor"
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
	ss.DB.Find(&sources)
	ss.Memstore.Set("sources", sources)
	return sources
}
