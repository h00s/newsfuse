package main

import (
	"github.com/go-raptor/raptor"
	"github.com/h00s/newsfuse/config"
	"github.com/h00s/newsfuse/config/initializers"
)

func main() {
	r := raptor.NewRaptor()

	r.Init(initializers.App(r.Utils.Config))
	r.RegisterRoutes(config.Routes())

	r.Listen()
}
