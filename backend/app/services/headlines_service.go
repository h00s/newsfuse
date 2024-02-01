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
	Scrapers        []internal.Scraper
	HeadlineChannel chan models.Headline
}

func NewHeadlinesService() *HeadlinesService {
	headlineChannel := make(chan models.Headline)

	hs := &HeadlinesService{
		Scrapers: []internal.Scraper{
			scrapers.NewKliknihr(headlineChannel),
			scrapers.NewMojportalhr(headlineChannel),
			scrapers.NewRadioDaruvar(headlineChannel),
		},
		HeadlineChannel: headlineChannel,
	}

	hs.OnInit(func() {
		hs.DB.AutoMigrate(&models.Headline{})

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
	hs.DB.Limit(50).Order("id desc").Find(&headlines)
	return headlines
}
