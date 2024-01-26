package services

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
	"github.com/h00s/newsfuse/internal/scrapers"
)

type HeadlinesService struct {
	raptor.Service
	Scrapers        []internal.Scraper
	Headlines       models.Headlines
	Headline        chan models.Headline
	storedHeadlines map[string]bool
}

func NewHeadlinesService() *HeadlinesService {
	headlineChan := make(chan models.Headline)

	return &HeadlinesService{
		Scrapers: []internal.Scraper{
			scrapers.NewKliknihr(headlineChan),
		},
		Headlines:       models.Headlines{},
		Headline:        headlineChan,
		storedHeadlines: make(map[string]bool),
	}
}

func (hs *HeadlinesService) Receive() {
	for {
		select {
		case h := <-hs.Headline:
			if _, ok := hs.storedHeadlines[h.URL]; !ok {
				hs.storedHeadlines[h.URL] = true
				hs.Headlines = append(hs.Headlines, h)
				hs.Utils.Log.Info("Received new headline", "Title", h.Title[:25])
			}
		}
	}
}
