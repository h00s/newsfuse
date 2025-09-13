package services

import (
	"context"
	"database/sql"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
	"github.com/uptrace/bun"
)

type StoriesService struct {
	raptor.Service
	Headlines *HeadlinesService

	claude *internal.Claude
}

func NewStoriesService() *StoriesService {
	ss := &StoriesService{}

	ss.OnInit(func() error {
		var err error
		ss.claude, err = internal.NewClaude(ss.Config.AppConfig["anthropic_key"])
		return err
	})

	return ss
}

func (ss *StoriesService) Get(headlineID int64) (*models.Story, error) {
	var story models.Story
	err := ss.Database.Conn().(*bun.DB).
		NewSelect().
		Model(&story).
		Where("headline_id = ?", headlineID).
		Scan(context.Background())

	if err != nil {
		if err == sql.ErrNoRows {
			story, err = ss.Scrape(headlineID)
			if err != nil {
				ss.Log.Error("Error scraping story", "error", err.Error())
				return nil, err
			}
			ss.Save(&story)
			return &story, nil
		}
		ss.Log.Error("Error getting story", "error", err.Error())
		return nil, err
	}

	return &story, nil
}

func (ss *StoriesService) Scrape(headlineID int64) (models.Story, error) {
	var headline models.Headline
	err := ss.Database.Conn().(*bun.DB).
		NewSelect().
		Model(&headline).
		Where("id = ?", headlineID).
		Scan(context.Background())

	if err != nil {
		return models.Story{}, err
	}

	content, err := ss.Headlines.Scrapers[int(headline.SourceID)].ScrapeStory(headline.URL)
	if err != nil {
		return models.Story{}, err
	}

	return models.Story{
		HeadlineID: headlineID,
		Content:    content,
	}, nil
}

func (ss *StoriesService) Save(story *models.Story) error {
	_, err := ss.Database.Conn().(*bun.DB).
		NewInsert().
		Model(story).
		Returning("id").
		Exec(context.Background())

	if err != nil {
		ss.Log.Error("Error creating story", "error", err.Error())
		return err
	}

	return nil
}

func (ss *StoriesService) Summarize(storyID int64) (models.Story, error) {
	var story models.Story
	err := ss.Database.Conn().(*bun.DB).
		NewSelect().
		Model(&story).
		Where("id = ?", storyID).
		Scan(context.Background())

	if err != nil {
		if err == sql.ErrNoRows {
			ss.Log.Error("Story not found", "storyID", storyID)
		} else {
			ss.Log.Error("Error getting story", "error", err.Error())
		}
		return story, err
	}

	if story.Summary != "" {
		return story, nil
	}

	summary, err := ss.claude.Summarize(story.Content)
	if err != nil {
		ss.Log.Error("Error summarizing story", "error", err.Error())
		return story, err
	}

	story.Summary = summary
	go ss.Database.Conn().(*bun.DB).
		NewUpdate().
		Model(&story).
		Column("summary").
		Where("id = ?", storyID).
		Exec(context.Background())

	return story, nil
}
