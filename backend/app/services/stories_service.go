package services

import (
	"fmt"

	"github.com/go-raptor/raptor/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal"
)

type StoriesService struct {
	raptor.Service
	Headlines *HeadlinesService

	claude *internal.Claude
}

func NewStoriesService() *StoriesService {
	ss := &StoriesService{}

	ss.OnInit(func() {
		var err error
		if ss.claude, err = internal.NewClaude(ss.Config.AppConfig["anthropic_key"].(string)); err != nil {
			ss.Log.Error("Error creating Claude", "error", err.Error())
		}
	})

	return ss
}

func (ss *StoriesService) Get(headlineID int) (*models.Story, error) {
	var story models.Story
	result := ss.DB.Where("headline_id = ?", headlineID).First(&story)
	if result.RowsAffected == 0 {
		story, err := ss.Scrape(headlineID)
		if err != nil {
			return nil, err
		}
		ss.Save(&story)
		return &story, nil
	}

	return &story, nil
}

func (ss *StoriesService) Scrape(headlineID int) (models.Story, error) {
	var headline models.Headline
	ss.DB.First(&headline, headlineID)
	content, err := ss.Headlines.Scrapers[int(headline.SourceID)].ScrapeStory(headline.URL)
	if err != nil {
		return models.Story{}, err
	}
	return models.Story{
		HeadlineID: uint(headlineID),
		Content:    content,
	}, nil
}

func (ss *StoriesService) Save(story *models.Story) error {
	result := ss.DB.Create(story)
	if result.Error != nil {
		ss.Log.Error("Error creating story", "DB", result.Error.Error())
		return result.Error
	}
	return nil
}

func (ss *StoriesService) Summarize(storyID int) (models.Story, error) {
	var story models.Story
	result := ss.DB.First(&story, storyID)
	if result.RowsAffected == 0 {
		return story, result.Error
	}
	if story.Summary != "" {
		return story, nil
	}
	summary, err := ss.claude.Summarize(story.Content)
	if err != nil {
		fmt.Println(err)
		return story, err
	}
	story.Summary = summary
	go ss.DB.Save(&story)
	return story, nil
}
