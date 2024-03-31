package services

import (
	"encoding/json"
	"fmt"
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
	Sources          *SourcesService
	Memstore         *Memstore
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
		result := hs.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&headlines)
		if result.Error == nil {
			hs.Log.Info("New headlines", "Count", result.RowsAffected)
			if result.RowsAffected > 0 {
				source := hs.Sources.Get(headlines[0].SourceID)
				if err := hs.allFromDB(int(source.TopicID), &headlines); err == nil {
					hs.Memstore.Set(fmt.Sprintf("headlines:%d", source.TopicID), headlines)
				}
			}
		} else {
			hs.Log.Error("Error creating headlines", "DB", result.Error.Error())
		}
	}
}

func (hs *HeadlinesService) All(topicID int) models.Headlines {
	var headlines models.Headlines
	data, err := hs.Memstore.Get(fmt.Sprintf("headlines:%d", topicID))
	if err == nil && data != "" {
		json.Unmarshal([]byte(data), &headlines)
		return headlines
	}

	if err := hs.allFromDB(topicID, &headlines); err != nil {
		hs.Log.Error("Error getting headlines", "Error", err.Error())
		return headlines
	}
	go hs.Memstore.Set(fmt.Sprintf("headlines:%d", topicID), headlines)

	return headlines
}

func (hs *HeadlinesService) allFromDB(topicID int, headlines *models.Headlines) error {
	return hs.DB.
		Limit(30).
		Order("id desc").
		Joins("JOIN sources ON headlines.source_id = sources.id").
		Where("sources.topic_id = ?", topicID).
		Find(&headlines).Error
}

func (hs *HeadlinesService) AllByLastID(topicID, lastID int) models.Headlines {
	var headlines models.Headlines
	if err := hs.DB.
		Limit(30).
		Order("id desc").
		Joins("JOIN sources ON headlines.source_id = sources.id").
		Where("sources.topic_id = ? AND headlines.id < ?", topicID, lastID).
		Find(&headlines).Error; err != nil {
		hs.Log.Error("Error getting headlines", "Error", err.Error())
	}

	return headlines
}
