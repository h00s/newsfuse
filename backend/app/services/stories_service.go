package services

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/models"
)

type StoriesService struct {
	raptor.Service
}

func (ss *StoriesService) Get(headlineID int) (*models.Story, error) {
	var story models.Story
	result := ss.DB.Where("headline_id = ?", headlineID).First(&story)
	if result.RowsAffected == 0 {
		story, err := ss.Scrape(headlineID)
		if err != nil {
			return nil, err
		}
		go ss.Save(&story)
	}

	return &story, nil
}

func (ss *StoriesService) Scrape(headlineID int) (models.Story, error) {
	headlines := ss.Services["HeadlinesService"].(*HeadlinesService)
	var headline models.Headline
	ss.DB.First(&headline, headlineID)
	content, err := headlines.Scrapers[int(headline.SourceID)].ScrapeStory(headline.URL)
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
