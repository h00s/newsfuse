package config

import "github.com/go-raptor/raptor"

func Routes() raptor.Routes {
	return raptor.CollectRoutes(
		raptor.Scope("api/v1",
			raptor.Route("GET", "/topics", "TopicsController", "All"),
			raptor.Route("GET", "/topics/:id/headlines", "HeadlinesController", "All"),

			raptor.Route("GET", "/headlines/:id/story", "HeadlinesController", "Story"),

			raptor.Route("GET", "/sources", "SourcesController", "All"),
		),

		raptor.Route("GET", "*", "SPAController", "Index"),
	)
}
