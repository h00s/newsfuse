package services

import (
	"slices"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/newsfuse/app/models"
)

// Topics and sources change only through migrations, so they stay cached for a long time.
const referenceDataTTL = time.Hour

type TopicsService struct {
	raptor.Service

	DB    *DatabaseService
	Cache *CacheService
}

func (s *TopicsService) List() (models.Topics, error) {
	return cached(s.Cache, "topics", referenceDataTTL, func() (models.Topics, error) {
		var topics models.Topics
		err := s.DB.Conn().NewSelect().
			Model(&topics).
			Order("topics.id").
			Scan(s.DB.Ctx)
		return topics, s.DB.HandleError(err)
	})
}

// Verify is what lists filtered by topic call first, so an unknown topic is a 404 rather than
// an empty list.
func (s *TopicsService) Verify(id int64) error {
	topics, err := s.List()
	if err != nil {
		return err
	}
	if !slices.ContainsFunc(topics, func(t models.Topic) bool { return t.ID == id }) {
		return errs.NewErrorNotFound("Topic not found")
	}
	return nil
}
