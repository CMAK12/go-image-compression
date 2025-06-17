package http

import (
	"context"
	"go-image-compression/internal/model"
	"go-image-compression/internal/service"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type (
	ImageService interface {
		Get(ctx context.Context, filter model.ListImageFilter) (multipart.File, error)
		Create(ctx context.Context, body []byte) error
	}

	Handler struct {
		ImageService ImageService
		logger       *zap.Logger
	}
)

func NewHandler(services service.Services, logger *zap.Logger) Handler {
	return Handler{
		ImageService: &services.ImageService,
		logger:       logger,
	}
}

func (h *Handler) SetupRoutes(app *fiber.App) {
	app.Use(loggerMiddleware(h.logger))

	api := app.Group("/api/v1")

	image := api.Group("/image")
	image.Get("", responseWrapper(h.getImage))
	image.Post("", h.createImage)
}
