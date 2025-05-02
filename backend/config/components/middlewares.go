package components

import (
	"github.com/go-raptor/middlewares/cors"
	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/core"
)

func Middlewares() raptor.Middlewares {
	return raptor.Middlewares{
		core.Use(cors.NewCORSMiddleware(cors.CORSConfig{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Authorization", "Content-Type"},
			AllowCredentials: true,
			MaxAge:           3600, // 1 hour
		})),
	}
}
