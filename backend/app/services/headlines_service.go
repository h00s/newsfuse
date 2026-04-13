package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
	"github.com/h00s/newsfuse/internal/scrapers"
	"github.com/uptrace/bun"
)

type HeadlinesService struct {
	raptor.Service

	Scrapers         map[int]internal.Scraper
	HeadlinesChannel chan models.Headlines
	Sources          *SourcesService
	Cache            *CacheService
}

func (s *HeadlinesService) Setup() error {
	s.HeadlinesChannel = make(chan models.Headlines)
	s.Scrapers = map[int]internal.Scraper{
		1:  scrapers.NewKliknihr(s.HeadlinesChannel, s.Log, 1),
		2:  scrapers.NewMojportalhr(s.HeadlinesChannel, s.Log, 2),
		3:  scrapers.NewRadioDaruvar(s.HeadlinesChannel, s.Log, 3),
		4:  scrapers.NewIndexhrCroatia(s.HeadlinesChannel, s.Log, 4),
		5:  scrapers.NewN1InfoCroatia(s.HeadlinesChannel, s.Log, 5),
		6:  scrapers.NewIndexhrWorld(s.HeadlinesChannel, s.Log, 6),
		7:  scrapers.NewN1InfoWorld(s.HeadlinesChannel, s.Log, 7),
		8:  scrapers.NewHackerNews(s.HeadlinesChannel, s.Log, 8),
		9:  scrapers.NewBughr(s.HeadlinesChannel, s.Log, 9),
		10: scrapers.NewTelegram(s.HeadlinesChannel, s.Log, 10),
		11: scrapers.NewHCL(s.HeadlinesChannel, s.Log, 11),
	}

	go s.Receive()
	for _, scraper := range s.Scrapers {
		scraper.Start()
	}

	return nil
}

func (s *HeadlinesService) Receive() {
	for {
		headlines := <-s.HeadlinesChannel
		slices.Reverse(headlines)
		newHeadlines := false
		for _, headline := range headlines {
			exists, err := s.Database.Conn().(*bun.DB).
				NewSelect().
				Model(&headline).
				Where("url = ?", headline.URL).
				Exists(context.Background())
			if err != nil {
				s.Log.Error("Error checking headline existence", "error", err.Error())
				continue
			}

			if !exists {
				_, err = s.Database.Conn().(*bun.DB).
					NewInsert().
					Model(&headline).
					Exec(context.Background())
				if err != nil {
					s.Log.Error("Error creating headline", "DB", err.Error())
					continue
				}
				newHeadlines = true
			}
		}
		if newHeadlines {
			source := s.Sources.Get(headlines[0].SourceID)
			if err := s.allFromDB(source.TopicID, &headlines); err == nil {
				go s.memstoreSetHeadlinesByTopicID(source.TopicID, &headlines)
			}
		}
	}
}

func (s *HeadlinesService) All(topicID int64) (models.Headlines, error) {
	var headlines models.Headlines

	if err := s.memstoreGetHeadlinesByTopicID(topicID, &headlines); err == nil {
		return headlines, nil
	}

	if err := s.allFromDB(topicID, &headlines); err != nil {
		return headlines, errs.NewErrorInternal(err.Error())
	}

	go s.memstoreSetHeadlinesByTopicID(topicID, &headlines)

	return headlines, nil
}

func (s *HeadlinesService) allFromDB(topicID int64, headlines *models.Headlines) error {
	if err := s.Database.Conn().(*bun.DB).
		NewSelect().
		Model(headlines).
		Join("JOIN sources s ON headline.source_id = s.id").
		Where("s.topic_id = ?", topicID).
		Order("headline.id desc").
		Limit(30).
		Scan(context.Background()); err != nil {
		s.Log.Error("Error getting headlines", "error", err.Error())
		return err
	}

	return nil
}

func (s *HeadlinesService) AllByLastID(topicID, lastID int64) (models.Headlines, error) {
	var headlines models.Headlines
	if err := s.Database.Conn().(*bun.DB).
		NewSelect().
		Model(&headlines).
		Join("JOIN sources s ON headline.source_id = s.id").
		Where("s.topic_id = ?", topicID).
		Where("headline.id < ?", lastID).
		Order("headline.id DESC").
		Limit(30).
		Scan(context.Background()); err != nil {
		s.Log.Error("Error getting headlines", "error", err)
		return headlines, errs.NewErrorInternal(err.Error())
	}

	return headlines, nil
}

func (s *HeadlinesService) Search(query string) (models.Headlines, error) {
	var headlines models.Headlines
	if err := s.Database.Conn().(*bun.DB).
		NewSelect().
		Model(&headlines).
		Where("title ILIKE ?", "%"+query+"%").
		Order("id DESC").
		Limit(100).
		Scan(context.Background()); err != nil {
		s.Log.Error("Error searching headlines", "error", err.Error())
		return headlines, errs.NewErrorInternal(err.Error())
	}

	return headlines, nil
}

func (s *HeadlinesService) Count(topicID int64, since time.Time) (int, error) {
	count, err := s.Database.Conn().(*bun.DB).
		NewSelect().
		Model((*models.Headline)(nil)).
		Join("JOIN sources s ON headline.source_id = s.id").
		Where("s.topic_id = ? AND headline.published_at > ?", topicID, since).
		Count(context.Background())
	if err != nil {
		s.Log.Error("Error counting headlines", "error", err.Error())
		return 0, errs.NewErrorInternal(err.Error())
	}

	return count, nil
}

func (s *HeadlinesService) memstoreGetHeadlinesByTopicID(topicID int64, headlines *models.Headlines) error {
	if data, ok := s.Cache.Get(fmt.Sprintf("headlines:%d", topicID)); ok {
		json.Unmarshal(data, headlines)
		return nil
	}
	s.Log.Warn("Headlines not found in memstore", "topic", topicID)
	return errors.New("headlines not found in memstore")
}

func (s *HeadlinesService) memstoreSetHeadlinesByTopicID(topicID int64, headlines *models.Headlines) {
	data, err := json.Marshal(headlines)
	if err != nil {
		s.Log.Warn("Error setting headlines in memstore", "topic", topicID, "error", err.Error())
	}
	s.Cache.Set(fmt.Sprintf("headlines:%d", topicID), data)
}
