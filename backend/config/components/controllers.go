package components

import (
	"github.com/go-raptor/controllers/spa"
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/app/controllers"
)

func Controllers() raptor.Controllers {
	return raptor.Controllers{
		&controllers.HeadlinesController{},
		&controllers.SourcesController{},
		&controllers.StoriesController{},
		&controllers.TopicsController{},
		spa.NewSPAController("public", "index.html"),
	}
}
