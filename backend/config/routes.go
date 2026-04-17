package config

import (
	_ "embed"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/router"
)

//go:embed routes.yaml
var routesYAML []byte

var routes = raptor.Must(router.ParseRoutesYAML(routesYAML))

func Routes() router.Routes {
	return routes
}
