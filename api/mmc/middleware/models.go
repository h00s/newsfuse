package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/h00s/newsfuse/api/mmc/models"
	"github.com/h00s/newsfuse/api/services"
)

type ModelsMiddleware struct {
	Sources *models.Sources
}

func NewModelsMiddleware(services *services.Services) *ModelsMiddleware {
	return &ModelsMiddleware{
		Sources: models.NewSources(services),
	}
}

func (m *ModelsMiddleware) ModelsMiddleware(c *fiber.Ctx) error {
	c.Locals("models", m)
	return c.Next()
}

func GetModels(c *fiber.Ctx) *ModelsMiddleware {
	return c.Locals("models").(*ModelsMiddleware)
}
