package config

import (
	"github.com/go-raptor/raptor/v3/router"
)

func Routes() router.Routes {
	return router.CollectRoutes(
		router.LoadFromYAML("config/routes.yaml"),
	)
}
