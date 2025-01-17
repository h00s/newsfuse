package config

import "github.com/go-raptor/raptor/v3"

func Routes() raptor.Routes {
	return raptor.CollectRoutes(
		raptor.Scope("api/v1",
			raptor.Scope("topics",
				raptor.Get("", "Topics#All"),
				raptor.Get(":id/headlines", "Headlines#All"),
				raptor.Get(":id/headlines/count", "Headlines#Count"),
			),

			raptor.Scope("headlines",
				raptor.Get(":id/story", "Stories#Get"),
				raptor.Get("search", "Headlines#Search"),
			),

			raptor.Scope("sources",
				raptor.Get("", "Sources#All"),
			),

			raptor.Scope("stories",
				raptor.Get(":id/summarize", "Stories#Summarize"),
			),
		),
	)
}
