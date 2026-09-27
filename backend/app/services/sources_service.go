package services

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/newsfuse/app/models"
)

type SourcesService struct {
	raptor.Service

	DB     *DatabaseService
	Cache  *CacheService
	Topics *TopicsService
}

func (s *SourcesService) List() (models.Sources, error) {
	return cached(s.Cache, "sources", referenceDataTTL, func() (models.Sources, error) {
		var sources models.Sources
		err := s.DB.Conn().NewSelect().
			Model(&sources).
			Order("sources.id").
			Scan(s.DB.Ctx)
		return sources, s.DB.HandleError(err)
	})
}

// ListByTopic verifies the topic first, so an unknown one is a 404 rather than an empty list.
func (s *SourcesService) ListByTopic(topicID int64) (models.Sources, error) {
	if err := s.Topics.Verify(topicID); err != nil {
		return nil, err
	}
	sources, err := s.List()
	if err != nil {
		return nil, err
	}
	var inTopic models.Sources
	for _, source := range sources {
		if source.TopicID == topicID {
			inTopic = append(inTopic, source)
		}
	}
	return inTopic, nil
}

func (s *SourcesService) Get(id int64) (*models.Source, error) {
	sources, err := s.List()
	if err != nil {
		return nil, err
	}
	for i := range sources {
		if sources[i].ID == id {
			source := sources[i] // a copy: the cached list is shared
			return &source, nil
		}
	}
	return nil, errs.NewErrorNotFound("Source not found")
}

// VerifyInTopic answers 404 for a source that doesn't exist or belongs to another topic.
func (s *SourcesService) VerifyInTopic(sourceID, topicID int64) error {
	source, err := s.Get(sourceID)
	if err != nil {
		return err
	}
	if source.TopicID != topicID {
		return errs.NewErrorNotFound("Source not found")
	}
	return nil
}
