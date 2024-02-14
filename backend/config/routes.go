package config

import "github.com/go-raptor/raptor"

func Routes() raptor.Routes {
	return raptor.CollectRoutes(
		raptor.Scope("api/v1",
			raptor.Scope("/topics",
				raptor.Route("GET", "", "TopicsController", "All"),
				raptor.Route("GET", "/:id/headlines", "HeadlinesController", "All"),
			),

			raptor.Scope("/headlines",
				raptor.Route("GET", "/:id/story", "StoriesController", "Get"),
			),

			raptor.Scope("/sources",
				raptor.Route("GET", "", "SourcesController", "All"),
			),
		),
		raptor.Route("GET", "*", "SPAController", "Index"),
	)
}
