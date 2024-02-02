package services

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/models"
)

type SourcesService struct {
	raptor.Service
}

func (ss *SourcesService) All() models.Sources {
	var sources models.Sources
	ss.DB.Find(&sources)
	return sources
}
