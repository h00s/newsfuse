package initializers

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/controllers"
	"github.com/h00s/newsfuse/app/services"
	"github.com/h00s/newsfuse/db"
)

func App() *raptor.AppInitializer {
	hs := services.NewHeadlinesService()
	ss := &services.SourcesService{}

	return &raptor.AppInitializer{
		Database: db.Migrations(),

		Services: raptor.Services{
			hs,
			ss,
		},

		Middlewares: raptor.Middlewares{},

		Controllers: raptor.Controllers{
			&controllers.HeadlinesController{
				Hs: hs,
			},
			&controllers.SourcesController{
				Ss: ss,
			},
			&controllers.SPAController{},
		},
	}
}
