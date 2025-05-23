package http

import (
	"go-image-compression/internal/service"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	ImageService service.ImageService
}

func NewHandler(services service.Services) Handler {
	return Handler{
		ImageService: services.ImageService,
	}
}

func (h *Handler) SetupRoutes(app *fiber.App) {
	app.Use(LoggerMiddleware())

	api := app.Group("/api/v1")

	image := api.Group("/image")
	image.Get("", ResponseWrapper(h.getImage))
	image.Post("", h.createImage)
}
