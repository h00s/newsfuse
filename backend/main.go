package main

import (
	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/newsfuse/config"
	"github.com/h00s/newsfuse/config/components"
)

func main() {
	app := raptor.New()

	app.Configure(components.New(app.Utils.Config))
	app.RegisterRoutes(config.Routes())
	app.Run()
}
