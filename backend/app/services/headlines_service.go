package services

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
	"github.com/h00s/newsfuse/internal/scrapers"
)

type HeadlinesService struct {
	raptor.Service
	Scrapers  []internal.Scraper
	Headlines models.Headlines
}

func NewHeadlinesService() *HeadlinesService {
	return &HeadlinesService{
		Scrapers: []internal.Scraper{
			scrapers.NewKliknihr(),
		},
	}
}
