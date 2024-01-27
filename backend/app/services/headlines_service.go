package services

import (
	"errors"

	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
	"github.com/h00s/newsfuse/internal/scrapers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type HeadlinesService struct {
	raptor.Service
	Scrapers []internal.Scraper
	Headline chan models.Headline
	db       *gorm.DB
}

func NewHeadlinesService() *HeadlinesService {
	db, err := gorm.Open(sqlite.Open("../newsfuse.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	db.AutoMigrate(&models.Headline{})
	headlineChan := make(chan models.Headline)

	return &HeadlinesService{
		Scrapers: []internal.Scraper{
			scrapers.NewKliknihr(headlineChan),
		},
		Headline: headlineChan,
		db:       db,
	}
}

func (hs *HeadlinesService) Receive() {
	for {
		select {
		case h := <-hs.Headline:
			result := hs.db.First(&h, "url = ?", h.URL)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				hs.db.Create(&h)
				hs.Utils.Log.Info("Added new headline", "Title", h.Title[:25])
			}
		}
	}
}

func (hs *HeadlinesService) All() models.Headlines {
	var headlines models.Headlines
	hs.db.Order("published_at desc").Find(&headlines)
	hs.Utils.Log.Info("Found headlines", "Count", len(headlines))
	return headlines
}
