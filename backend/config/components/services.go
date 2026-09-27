package components

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/app/services"
)

// Services are set up in this order, and cleaned up in reverse: the scrapers stop before the
// database and cache they write to.
func Services() raptor.Services {
	return raptor.Services{
		&services.DatabaseService{},
		&services.CacheService{},
		&services.TopicsService{},
		&services.SourcesService{},
		&services.HeadlinesService{},
		&services.GenAIService{},
		&services.ScrapersService{},
		&services.StoriesService{},
	}
}
