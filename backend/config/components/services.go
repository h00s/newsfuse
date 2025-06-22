package components

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/app/services"
)

func Services(c *raptor.Config) raptor.Services {
	return raptor.Services{
		services.NewHeadlinesService(),
		&services.SourcesService{},
		services.NewStoriesService(),
		&services.TopicsService{},
		services.NewCacheService(),
	}
}
