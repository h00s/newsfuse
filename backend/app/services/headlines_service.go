package services

import (
	"errors"

	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
	"github.com/h00s/newsfuse/internal/scrapers"
	"gorm.io/gorm"
)

type HeadlinesService struct {
	raptor.Service
	Scrapers        map[int]internal.Scraper
	HeadlineChannel chan models.Headline
}

func NewHeadlinesService() *HeadlinesService {
	headlineChannel := make(chan models.Headline)

	hs := &HeadlinesService{
		Scrapers: map[int]internal.Scraper{
			1: scrapers.NewKliknihr(headlineChannel, 1),
			2: scrapers.NewMojportalhr(headlineChannel, 2),
			3: scrapers.NewRadioDaruvar(headlineChannel, 3),
		},
		HeadlineChannel: headlineChannel,
	}

	hs.OnInit(func() {
		go hs.Receive()
		for _, s := range hs.Scrapers {
			s.Init(hs.Utils)
			s.Start()
		}
	})

	return hs
}

func (hs *HeadlinesService) Receive() {
	for {
		select {
		case h := <-hs.HeadlineChannel:
			result := hs.DB.First(&h, "url = ?", h.URL)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				hs.DB.Create(&h)
				hs.Log.Info("Added new headline", "Title", h.Title)
			}
		}
	}
}

func (hs *HeadlinesService) All() models.Headlines {
	var headlines models.Headlines
	hs.DB.Limit(50).Order("id desc").Preload("Source").Find(&headlines)
	return headlines
}

func (hs *HeadlinesService) Story(id int) (models.Headline, error) {
	hs.Log.Info("Story", "ID", id)
	var headline models.Headline
	return headline, nil
}
