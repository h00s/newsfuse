package initializers

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/controllers"
	"github.com/h00s/newsfuse/app/services"
)

func App() *raptor.AppInitializer {
	hs := services.NewHeadlinesService()

	return &raptor.AppInitializer{
		Services: raptor.Services{
			hs,
		},

		Middlewares: raptor.Middlewares{},

		Controllers: raptor.Controllers{
			&controllers.HeadlinesController{
				Hs: hs,
			},
			&controllers.SPAController{},
		},
	}
}
