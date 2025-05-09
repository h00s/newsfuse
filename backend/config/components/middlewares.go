package components

import (
	"github.com/go-raptor/middlewares/cors"
	"github.com/go-raptor/middlewares/logger"
	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/core"
)

func Middlewares(c *raptor.Config) raptor.Middlewares {
	return raptor.Middlewares{
		core.Use(cors.NewCORSMiddleware(cors.CORSConfig{
			AllowOrigins:     []string{c.AppConfig["cors_allow_origins"]},
			AllowMethods:     cors.DefaultCORSConfig.AllowMethods,
			AllowHeaders:     cors.DefaultCORSConfig.AllowHeaders,
			AllowCredentials: true,
			MaxAge:           3600,
		})),
		core.Use(&logger.LoggerMiddleware{}),
	}
}
