package services

import (
	"context"
	"database/sql"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/app/models"
	"github.com/uptrace/bun"
)

type StoriesService struct {
	raptor.Service

	Headlines *HeadlinesService
	GenAI     *GenAIService
}

func (s *StoriesService) Get(headlineID int64) (*models.Story, error) {
	var story models.Story
	err := s.Database.Conn().(*bun.DB).
		NewSelect().
		Model(&story).
		Where("headline_id = ?", headlineID).
		Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			story, err = s.Scrape(headlineID)
			if err != nil {
				s.Log.Error("Error scraping story", "error", err.Error())
				return nil, err
			}
			s.Save(&story)
			return &story, nil
		}
		s.Log.Error("Error getting story", "error", err.Error())
		return nil, err
	}

	return &story, nil
}

func (s *StoriesService) Scrape(headlineID int64) (models.Story, error) {
	var headline models.Headline
	err := s.Database.Conn().(*bun.DB).
		NewSelect().
		Model(&headline).
		Where("id = ?", headlineID).
		Scan(context.Background())
	if err != nil {
		return models.Story{}, err
	}

	content, err := s.Headlines.Scrapers[int(headline.SourceID)].ScrapeStory(headline.URL)
	if err != nil {
		return models.Story{}, err
	}

	return models.Story{
		HeadlineID: headlineID,
		Content:    content,
	}, nil
}

func (s *StoriesService) Save(story *models.Story) error {
	_, err := s.Database.Conn().(*bun.DB).
		NewInsert().
		Model(story).
		Returning("id").
		Exec(context.Background())
	if err != nil {
		s.Log.Error("Error creating story", "error", err.Error())
		return err
	}

	return nil
}

func (s *StoriesService) Summarize(storyID int64) (models.Story, error) {
	var story models.Story
	err := s.Database.Conn().(*bun.DB).
		NewSelect().
		Model(&story).
		Where("id = ?", storyID).
		Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			s.Log.Error("Story not found", "storyID", storyID)
		} else {
			s.Log.Error("Error getting story", "error", err.Error())
		}
		return story, err
	}

	if story.Summary != "" {
		return story, nil
	}

	summary, err := s.GenAI.Summarize(story.Content)
	if err != nil {
		s.Log.Error("Error summarizing story", "error", err.Error())
		return story, err
	}

	story.Summary = summary
	go s.Database.Conn().(*bun.DB).
		NewUpdate().
		Model(&story).
		Column("summary").
		Where("id = ?", storyID).
		Exec(context.Background())

	return story, nil
}
