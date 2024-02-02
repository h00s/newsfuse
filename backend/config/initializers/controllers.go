package initializers

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/app/controllers"
	"github.com/h00s/newsfuse/app/services"
)

func ServicesAndControllers() (raptor.Services, raptor.Controllers) {
	hs := services.NewHeadlinesService()
	ss := &services.SourcesService{}

	return raptor.Services{
			hs,
			ss,
		}, raptor.Controllers{
			&controllers.HeadlinesController{
				Hs: hs,
			},
			&controllers.SourcesController{
				Ss: ss,
			},
			&controllers.SPAController{},
		}
}
