package main

import (
	"github.com/go-raptor/raptor/v3"
	"github.com/h00s/newsfuse/app/utils"
	"github.com/h00s/newsfuse/config"
	"github.com/h00s/newsfuse/config/components"
)

func main() {
	app := raptor.New()

	logistiq, err := utils.NewLogistiqHandler(app.Utils)
	if err != nil {
		app.Utils.Log.Error("Failed to create Logistiq handler", "error", err)
		return
	}
	defer logistiq.Close()

	app.Configure(components.New(app.Utils.Config))
	app.RegisterRoutes(config.Routes())
	app.Run()
}
