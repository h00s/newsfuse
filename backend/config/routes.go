package config

import "github.com/go-raptor/raptor"

func Routes() raptor.Routes {
	return raptor.CollectRoutes(
		raptor.Route("GET", "/api/v1/headlines", "HeadlinesController", "All"),
		raptor.Route("GET", "/api/v1/sources", "SourcesController", "All"),
		raptor.Route("GET", "*", "SPAController", "Index"),
	)
}
