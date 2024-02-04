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
		h := <-hs.HeadlineChannel
		result := hs.DB.First(&h, "url = ?", h.URL)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			result := hs.DB.Create(&h)
			if result.Error != nil {
				hs.Log.Error("Error creating headline", "DB", result.Error.Error())
			}
		}
	}
}

func (hs *HeadlinesService) All() models.Headlines {
	var headlines models.Headlines
	hs.DB.Limit(50).Order("id desc").Preload("Source").Find(&headlines)
	return headlines
}

func (hs *HeadlinesService) GetStory(headlineID int) (*models.Story, error) {
	var headline models.Headline
	var story models.Story
	result := hs.DB.Where("headline_id = ?", headlineID).First(&story)
	if result.RowsAffected == 0 {
		hs.DB.First(&headline, headlineID)
		content, err := hs.Scrapers[int(headline.SourceID)].ScrapeStory(headline.URL)
		if err != nil {
			return &story, err
		}
		story = models.Story{
			HeadlineID: uint(headlineID),
			Content:    content,
		}
		result := hs.DB.Create(&story)
		if result.Error != nil {
			hs.Log.Error("Error creating story", "Headline", headlineID)
		}
	}

	return &story, nil
}
