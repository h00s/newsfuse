package components

import (
	"github.com/go-raptor/middlewares/cors"
	"github.com/go-raptor/middlewares/logger"
	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/core"
)

func Middlewares() raptor.Middlewares {
	return raptor.Middlewares{
		core.Use(&logger.LoggerMiddleware{}),
		core.Use(&cors.CORSMiddleware{}),
	}
}
