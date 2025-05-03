package http

import "github.com/gofiber/fiber/v2"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1", LoggerMiddleware())

	compress := api.Group("/compress")
	compress.Get("", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
}
