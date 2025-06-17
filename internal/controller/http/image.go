package http

import (
	"fmt"
	"mime/multipart"

	"go-image-compression/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/nordew/go-errx"
	"go.uber.org/zap"
)

func (h *Handler) getImage(c *fiber.Ctx) (multipart.File, error) {
	var filter model.ListImageFilter

	if err := c.QueryParser(&filter); err != nil {
		return nil, errx.NewBadRequest().WithDescription("failed to parse query params")
	}
	if filter.ID == "" {
		return nil, errx.NewBadRequest().WithDescription("image ID is required")
	}

	image, err := h.ImageService.Get(c.Context(), filter)
	if err != nil {
		return nil, fmt.Errorf("handler.image.getImage: %w", err)
	}

	return image, nil
}

func (h *Handler) createImage(c *fiber.Ctx) error {
	body := c.Body()

	if err := h.ImageService.Create(c.Context(), body); err != nil {
		h.logger.Error("handler.image.createImage: %v", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create image")
	}

	c.Status(fiber.StatusCreated)
	return nil
}
