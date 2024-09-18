package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
	"github.com/h00s/newsfuse/internal/scrapers"
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
			1:  scrapers.NewKliknihr(headlinesChannel, 1),
			2:  scrapers.NewMojportalhr(headlinesChannel, 2),
			3:  scrapers.NewRadioDaruvar(headlinesChannel, 3),
			4:  scrapers.NewIndexhrCroatia(headlinesChannel, 4),
			5:  scrapers.NewN1InfoCroatia(headlinesChannel, 5),
			6:  scrapers.NewIndexhrWorld(headlinesChannel, 6),
			7:  scrapers.NewN1InfoWorld(headlinesChannel, 7),
			8:  scrapers.NewHackerNews(headlinesChannel, 8),
			9:  scrapers.NewBughr(headlinesChannel, 9),
			10: scrapers.NewTelegram(headlinesChannel, 10),
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
		newHeadlines := false
		for _, headline := range headlines {
			exists, err := hs.DB.
				NewSelect().
				Model(&headline).
				Where("url = ?", headline.URL).
				Exists(context.Background())
			if err != nil {
				hs.Log.Error("Error checking headline existence", "error", err.Error())
				continue
			}

			if !exists {
				_, err = hs.DB.
					NewInsert().
					Model(&headline).
					Exec(context.Background())
				if err != nil {
					hs.Log.Error("Error creating headline", "DB", err.Error())
					continue
				}
				newHeadlines = true
			}
		}
		if newHeadlines {
			source := hs.Sources.Get(headlines[0].SourceID)
			if err := hs.allFromDB(source.TopicID, &headlines); err == nil {
				go hs.memstoreSetHeadlinesByTopicID(source.TopicID, &headlines)
			}
		}
	}
}

func (hs *HeadlinesService) All(topicID int64) (models.Headlines, error) {
	var headlines models.Headlines

	if err := hs.memstoreGetHeadlinesByTopicID(topicID, &headlines); err == nil {
		return headlines, nil
	}

	if err := hs.allFromDB(topicID, &headlines); err != nil {
		return headlines, raptor.NewErrorInternal(err.Error())
	}

	go hs.memstoreSetHeadlinesByTopicID(topicID, &headlines)

	return headlines, nil
}

func (hs *HeadlinesService) allFromDB(topicID int64, headlines *models.Headlines) error {
	if err := hs.DB.
		NewSelect().
		Model(headlines).
		Join("JOIN sources s ON headline.source_id = s.id").
		Where("s.topic_id = ?", topicID).
		Order("headline.id desc").
		Limit(30).
		Scan(context.Background()); err != nil {
		hs.Log.Error("Error getting headlines", "error", err.Error())
		return err
	}

	return nil
}

func (hs *HeadlinesService) AllByLastID(topicID, lastID int64) (models.Headlines, error) {
	var headlines models.Headlines
	if err := hs.DB.
		NewSelect().
		Model(&headlines).
		Join("JOIN sources s ON headline.source_id = s.id").
		Where("s.topic_id = ?", topicID).
		Where("headline.id < ?", lastID).
		Order("headline.id DESC").
		Limit(30).
		Scan(context.Background()); err != nil {
		hs.Log.Error("Error getting headlines", "error", err)
		return headlines, raptor.NewErrorInternal(err.Error())
	}

	return headlines, nil
}

func (hs *HeadlinesService) Search(query string) (models.Headlines, error) {
	var headlines models.Headlines
	if err := hs.DB.
		NewSelect().
		Model(&headlines).
		Where("title ILIKE ?", "%"+query+"%").
		Order("id DESC").
		Limit(100).
		Scan(context.Background()); err != nil {
		hs.Log.Error("Error searching headlines", "error", err.Error())
		return headlines, raptor.NewErrorInternal(err.Error())
	}

	return headlines, nil
}

func (hs *HeadlinesService) Count(topicID int64, since time.Time) (int, error) {
	count, err := hs.DB.
		NewSelect().
		Model((*models.Headline)(nil)).
		Join("JOIN sources s ON headline.source_id = s.id").
		Where("s.topic_id = ? AND headline.published_at > ?", topicID, since).
		Count(context.Background())

	if err != nil {
		hs.Log.Error("Error counting headlines", "error", err.Error())
		return 0, raptor.NewErrorInternal(err.Error())
	}

	return count, nil
}

func (hs *HeadlinesService) memstoreGetHeadlinesByTopicID(topicID int64, headlines *models.Headlines) error {
	data, err := hs.Memstore.Get(fmt.Sprintf("headlines:%d", topicID))
	if err == nil && data != "" {
		json.Unmarshal([]byte(data), headlines)
		return nil
	}

	hs.Log.Warn("Headlines not found in memstore", "topic", topicID)
	return errors.New("headlines not found in memstore")
}

func (hs *HeadlinesService) memstoreSetHeadlinesByTopicID(topicID int64, headlines *models.Headlines) {
	err := hs.Memstore.Set(fmt.Sprintf("headlines:%d", topicID), headlines)
	if err != nil {
		hs.Log.Warn("Error setting headlines in memstore", "topic", topicID, "error", err.Error())
	}
}
