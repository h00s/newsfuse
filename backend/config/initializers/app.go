package initializers

import (
	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/newsfuse/config"
)

func App(c *raptor.Config) *raptor.AppInitializer {
	return &raptor.AppInitializer{
		Routes:      config.Routes(),
		Database:    Database(),
		Services:    Services(c),
		Middlewares: Middlewares(),
		Controllers: Controllers(),
	}
}
