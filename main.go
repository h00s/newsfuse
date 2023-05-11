package main

import (
	"log"
	"os"

	"github.com/h00s/newsfuse/api"
	"github.com/h00s/newsfuse/api/config"
)

func main() {
	config := config.NewConfig()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	api := api.NewAPI(config, logger)
	api.Start()
	api.WaitForShutdown()
}
