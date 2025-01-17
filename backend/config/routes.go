package config

import "github.com/go-raptor/raptor/v3"

func Routes() raptor.Routes {
	return raptor.CollectRoutes(
		raptor.Scope("api/v1",
			raptor.Scope("topics",
				raptor.Get("", "TopicsController", "All"),
				raptor.Get(":id/headlines", "HeadlinesController", "All"),
				raptor.Get(":id/headlines/count", "HeadlinesController", "Count"),
			),

			raptor.Scope("headlines",
				raptor.Get(":id/story", "StoriesController", "Get"),
				raptor.Get("search", "HeadlinesController", "Search"),
			),

			raptor.Scope("sources",
				raptor.Get("", "SourcesController", "All"),
			),

			raptor.Scope("stories",
				raptor.Get(":id/summarize", "StoriesController", "Summarize"),
			),
		),
	)
}
