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
				raptor.Route("GET", "/:id/story", "HeadlinesController", "Story"),
			),

			raptor.Scope("/sources",
				raptor.Route("GET", "/sources", "SourcesController", "All"),
			),
		),
		raptor.Route("GET", "*", "SPAController", "Index"),
	)
}
