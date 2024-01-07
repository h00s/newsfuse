package initializers

import (
	"github.com/h00s/newsfuse/app/controllers"
	"github.com/h00s/raptor"
)

func App() *raptor.AppInitializer {
	return &raptor.AppInitializer{
		Services: raptor.Services{},

		Middlewares: raptor.Middlewares{},

		Controllers: raptor.Controllers{
			&controllers.SPAController{},
		},
	}
}
