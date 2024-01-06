package config

import (
	"github.com/h00s/newsfuse/app/controllers"
	"github.com/h00s/raptor"
)

func Controllers() raptor.Controllers {
	return raptor.Controllers{
		&controllers.SPAController{},
	}
}
