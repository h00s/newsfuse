package services

import (
	"slices"

	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
	"github.com/h00s/newsfuse/internal/scrapers"
	"gorm.io/gorm/clause"
)

type HeadlinesService struct {
	raptor.Service
	Scrapers         map[int]internal.Scraper
	HeadlinesChannel chan models.Headlines
}

func NewHeadlinesService() *HeadlinesService {
	headlinesChannel := make(chan models.Headlines)

	hs := &HeadlinesService{
		Scrapers: map[int]internal.Scraper{
			1: scrapers.NewKliknihr(headlinesChannel, 1),
			2: scrapers.NewMojportalhr(headlinesChannel, 2),
			3: scrapers.NewRadioDaruvar(headlinesChannel, 3),
			4: scrapers.NewIndexhrCroatia(headlinesChannel, 4),
			5: scrapers.NewN1InfoCroatia(headlinesChannel, 5),
			6: scrapers.NewIndexhrWorld(headlinesChannel, 6),
			7: scrapers.NewN1InfoWorld(headlinesChannel, 7),
			8: scrapers.NewHackerNews(headlinesChannel, 8),
		},
		HeadlinesChannel: headlinesChannel,
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
		headlines := <-hs.HeadlinesChannel
		slices.Reverse(headlines)
		if result := hs.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&headlines); result.Error != nil {
			hs.Log.Error("Error creating headlines", "DB", result.Error.Error())
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
