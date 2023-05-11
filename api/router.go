package api

import "github.com/gofiber/fiber/v2"

func (api *API) setRoutes() {
	api.server.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.Send([]byte("Hello, World!"))
	})
}
