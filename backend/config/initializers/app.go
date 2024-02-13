package initializers

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/db"
)

func App() *raptor.AppInitializer {
	return &raptor.AppInitializer{
		Database:    db.Migrations(),
		Services:    Services(),
		Middlewares: Middlewares(),
		Controllers: Controllers(),
	}
}
