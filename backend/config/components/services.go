package components

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/app/services"
)

func Services() raptor.Services {
	return raptor.Services{
		&services.SourcesService{},
		&services.StoriesService{},
		&services.TopicsService{},
		&services.CacheService{},
		&services.GenAIService{},
		&services.HeadlinesService{},
	}
}
