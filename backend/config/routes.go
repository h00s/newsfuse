package config

import (
	"github.com/go-raptor/raptor/v3/router"
)

func Routes() router.Routes {
	return router.CollectRoutes(
		router.Scope("api/v1",
			router.Scope("topics",
				router.Get("", "Topics.All"),
				router.Get(":id/headlines", "Headlines.All"),
				router.Get(":id/headlines/count", "Headlines.Count"),
			),

			router.Scope("headlines",
				router.Get(":id/story", "Stories.Get"),
				router.Get("search", "Headlines.Search"),
			),

			router.Scope("sources",
				router.Get("", "Sources.All"),
			),

			router.Scope("stories",
				router.Get(":id/summarize", "Stories.Summarize"),
			),
		),
	)
}
