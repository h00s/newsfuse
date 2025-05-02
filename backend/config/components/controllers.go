package components

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/app/controllers"
)

func Controllers() raptor.Controllers {
	return raptor.Controllers{
		&controllers.HeadlinesController{},
		&controllers.SourcesController{},
		&controllers.StoriesController{},
		&controllers.TopicsController{},
		controllers.NewSPAController("public", "index.html"),
	}
}
