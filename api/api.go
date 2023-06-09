package api

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/h00s/newsfuse/api/config"
	"github.com/h00s/newsfuse/api/mmc/middleware"
	"github.com/h00s/newsfuse/api/robot"
	"github.com/h00s/newsfuse/api/services"
)

type API struct {
	config   *config.Config
	server   *fiber.App
	services *services.Services
	robot    *robot.Robot
}

func NewAPI(config *config.Config, logger *log.Logger) *API {
	server := fiber.New()
	services := services.NewServices(config, logger)
	servicesMiddleware := middleware.NewServicesMiddleware(services)
	limiterMiddleware := middleware.NewLimiterMiddleware(&config.Limiter)
	modelsMiddleware := middleware.NewModelsMiddleware(services)

	server.Use(servicesMiddleware.ServicesMiddleware)
	server.Use(limiterMiddleware.LimiterMiddleware())
	server.Use(modelsMiddleware.ModelsMiddleware)

	return &API{
		config:   config,
		server:   server,
		services: services,
		robot:    robot.NewRobot(modelsMiddleware),
	}
}

func (api *API) Start() {
	api.services.Logger.Println("Starting server on :8080")
	api.setRoutes()
	go func() {
		if err := api.server.Listen(":8080"); err != nil && err != http.ErrServerClosed {
			api.services.Logger.Fatal("Error starting server: " + err.Error())
		}
	}()
	api.robot.Start()
}

func (api *API) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	api.services.Close()
	if err := api.server.Shutdown(); err != nil {
		api.services.Logger.Fatal(err)
	}
}
