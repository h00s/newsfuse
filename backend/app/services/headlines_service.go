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
	hs := &HeadlinesService{}
	hs.Scrapers = []internal.Scraper{
		scrapers.NewKliknihr(hs.AddHeadline),
	}
	hs.Headlines = models.Headlines{}
	return hs
}

func (hs *HeadlinesService) AddHeadline(h models.Headline) {
	hs.Headlines = append(hs.Headlines, h)
}
