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
			4: scrapers.NewIndexhrCroatia(headlineChannel, 4),
			5: scrapers.NewN1InfoCroatia(headlineChannel, 5),
			6: scrapers.NewIndexhrWorld(headlineChannel, 6),
			7: scrapers.NewN1InfoWorld(headlineChannel, 7),
			8: scrapers.NewHackerNews(headlineChannel, 8),
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

func (hs *HeadlinesService) All(topicID int) models.Headlines {
	var headlines models.Headlines
	hs.DB.
		Limit(50).
		Order("id desc").
		Joins("JOIN sources ON headlines.source_id = sources.id").
		Where("sources.topic_id = ?", topicID).
		Find(&headlines)
	return headlines
}
