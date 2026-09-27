package components

import (
	"github.com/go-raptor/controllers/spa/v2"
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/app/controllers"
)

func Controllers() raptor.Controllers {
	return raptor.Controllers{
		&controllers.HeadlinesController{},
		&controllers.SourcesController{},
		&controllers.StoriesController{},
		&controllers.TopicsController{},
		&spa.SPAController{}, // serves public/; spa_optional in the dev and test configs
	}
}
