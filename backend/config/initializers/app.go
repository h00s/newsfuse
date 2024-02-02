package initializers

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/db"
)

func App() *raptor.AppInitializer {
	services, controllers := ServicesAndControllers()

	return &raptor.AppInitializer{
		Database:    db.Migrations(),
		Services:    services,
		Middlewares: Middlewares(),
		Controllers: controllers,
	}
}
