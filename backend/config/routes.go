package config

import "github.com/go-raptor/raptor/v3"

func Routes() raptor.Routes {
	return raptor.CollectRoutes(
		raptor.Scope("api/v1",
			raptor.Scope("/topics",
				raptor.Route("GET", "", "TopicsController", "All"),
				raptor.Route("GET", "/:id/headlines", "HeadlinesController", "All"),
				raptor.Route("GET", "/:id/headlines/count", "HeadlinesController", "Count"),
			),

			raptor.Scope("/headlines",
				raptor.Route("GET", "/:id/story", "StoriesController", "Get"),
				raptor.Route("GET", "/search", "HeadlinesController", "Search"),
			),

			raptor.Scope("/sources",
				raptor.Route("GET", "", "SourcesController", "All"),
			),

			raptor.Scope("/stories",
				raptor.Route("GET", "/:id/summarize", "StoriesController", "Summarize"),
			),
		),
	)
}
