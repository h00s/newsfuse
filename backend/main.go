package main

import (
	"github.com/h00s/newsfuse/config"
	"github.com/h00s/newsfuse/config/initializers"
	"github.com/h00s/raptor"
)

func main() {
	r := raptor.NewRaptor()

	r.Init(initializers.App())
	r.Routes(config.Routes())

	r.Listen()
}
