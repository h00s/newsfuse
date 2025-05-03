package main

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/app/utils"
	"github.com/h00s/newsfuse/config"
	"github.com/h00s/newsfuse/config/components"
)

func main() {
	app := raptor.New()

	logistiq, err := utils.NewLogistiqHandler(app.Core.Resources.Config)
	if err != nil {
		app.Core.Resources.Log.Error("Failed to create Logistiq handler", "error", err)
		return
	}
	defer logistiq.Close()

	app.Configure(components.New(app.Core.Resources.Config))
	app.RegisterRoutes(config.Routes())
	app.Run()
}
