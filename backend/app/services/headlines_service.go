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
	Headline  chan models.Headline
}

func NewHeadlinesService() *HeadlinesService {
	headlineChan := make(chan models.Headline)

	return &HeadlinesService{
		Scrapers: []internal.Scraper{
			scrapers.NewKliknihr(headlineChan),
		},
		Headlines: models.Headlines{},
		Headline:  headlineChan,
	}
}

func (hs *HeadlinesService) Receive() {
	for {
		select {
		case h := <-hs.Headline:
			hs.Headlines = append(hs.Headlines, h)
		}
	}
}
