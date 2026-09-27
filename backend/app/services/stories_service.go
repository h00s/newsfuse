package services

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/app/utils"
	"golang.org/x/sync/singleflight"
)

type StoriesService struct {
	raptor.Service

	DB        *DatabaseService
	Headlines *HeadlinesService
	Scrapers  *ScrapersService
	GenAI     *GenAIService

	// Concurrent requests for the same story share one scrape and one LLM call.
	scraping    singleflight.Group
	summarizing singleflight.Group
}

func (s *StoriesService) Get(id int64) (*models.Story, error) {
	story := new(models.Story)
	err := s.DB.Conn().NewSelect().
		Model(story).
		Where("stories.id = ?", id).
		Scan(s.DB.Ctx)
	if err != nil {
		return nil, s.DB.HandleErrorNotFound(err, "Story not found")
	}
	story.Content = utils.SanitizeStory(story.Content)
	return story, nil
}

// GetByHeadline returns the story behind a headline, scraping and storing it on the first read.
// Content is sanitized again on every read, which also covers stories stored before scraping
// sanitized them.
func (s *StoriesService) GetByHeadline(ctx context.Context, headlineID int64) (*models.Story, error) {
	story, err := s.findByHeadline(headlineID)
	if err != nil || story != nil {
		return story, err
	}
	v, err, _ := s.scraping.Do(strconv.FormatInt(headlineID, 10), func() (any, error) {
		return s.scrape(ctx, headlineID)
	})
	if err != nil {
		return nil, err
	}
	return v.(*models.Story), nil
}

// findByHeadline returns nil without an error when the headline has no story yet.
func (s *StoriesService) findByHeadline(headlineID int64) (*models.Story, error) {
	story := new(models.Story)
	err := s.DB.Conn().NewSelect().
		Model(story).
		Where("stories.headline_id = ?", headlineID).
		Scan(s.DB.Ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, s.DB.HandleError(err)
	}
	story.Content = utils.SanitizeStory(story.Content)
	return story, nil
}

func (s *StoriesService) scrape(ctx context.Context, headlineID int64) (*models.Story, error) {
	headline, err := s.Headlines.Get(headlineID)
	if err != nil {
		return nil, err
	}
	if headline.Source == nil || !headline.Source.IsScrapable {
		return nil, errs.NewErrorNotFound("Story not available")
	}
	content, err := s.Scrapers.ScrapeStory(ctx, headline.SourceID, headline.URL)
	if err != nil {
		return nil, err
	}

	// Another instance may have stored the story meanwhile; its row wins.
	story := &models.Story{HeadlineID: headlineID, Content: content}
	if _, err := s.DB.Conn().NewInsert().
		Model(story).
		On("CONFLICT (headline_id) DO NOTHING").
		Exec(s.DB.Ctx); err != nil {
		return nil, s.DB.HandleError(err)
	}
	stored, err := s.findByHeadline(headlineID)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, errs.NewErrorNotFound("Story not found") // the headline was deleted meanwhile
	}
	return stored, nil
}

// Summarize summarizes a story once and stores the result; later calls return the stored one.
func (s *StoriesService) Summarize(ctx context.Context, id int64) (*models.Story, error) {
	story, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	if story.Summary != "" {
		return story, nil
	}
	v, err, _ := s.summarizing.Do(strconv.FormatInt(id, 10), func() (any, error) {
		return s.summarize(ctx, story)
	})
	if err != nil {
		return nil, err
	}
	return v.(*models.Story), nil
}

func (s *StoriesService) summarize(ctx context.Context, story *models.Story) (*models.Story, error) {
	summary, err := s.GenAI.Summarize(ctx, story.Content)
	if err != nil {
		s.Log.Error("Summarization failed", "story", story.ID, "error", err)
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, errs.NewErrorGatewayTimeout("Summarization timed out")
		}
		return nil, errs.NewErrorBadGateway("Summarization failed")
	}

	updated := &models.Story{ID: story.ID, Summary: summary}
	res, err := s.DB.Conn().NewUpdate().
		Model(updated).
		Column("summary").
		Where("stories.id = ?", story.ID).
		Returning("*").
		Exec(s.DB.Ctx)
	if err := s.DB.HandleAffected(res, err, "Story not found"); err != nil {
		return nil, err
	}
	updated.Content = utils.SanitizeStory(updated.Content)
	return updated, nil
}
