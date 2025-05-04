package http

import (
	"go-image-compression/internal/service"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Services service.Services
}

func NewHandler(services service.Services) Handler {
	return Handler{
		Services: services,
	}
}

func (h *Handler) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1", LoggerMiddleware())

	image := api.Group("/image")
	image.Get("", h.listImage)
	image.Post("", h.createImage)
}
