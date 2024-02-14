package initializers

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/services"
)

func Services() raptor.Services {
	return raptor.Services{
		services.NewHeadlinesService(),
		&services.SourcesService{},
		&services.StoriesService{},
		&services.TopicsService{},
	}
}
