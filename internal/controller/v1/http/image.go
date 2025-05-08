package http

import (
	"fmt"
	"io"
	"log"

	"go-image-compression/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/nordew/go-errx"
)

const codepath = "controller/v1/http/image.go"

func (h *Handler) getImage(c *fiber.Ctx) (io.Reader, error) {
	var filter model.ListImageFilter

	if err := c.QueryParser(&filter); err != nil {
		return nil, errx.NewBadRequest().WithDescription("failed to parse query params")
	}
	if filter.CompressPercent == 0 || filter.ID == "" {
		log.Printf("%s: missing required params: compress_percent=%d, id=%s", codepath, filter.CompressPercent, filter.ID)
		return nil, errx.NewBadRequest().WithDescription("compress_percent and id are required")
	}

	image, err := h.ImageService.Get(c.Context(), filter)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", codepath, err)
	}

	return image, nil
}

func (h *Handler) createImage(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Printf("%s: %v", codepath, err)
		return fiber.NewError(fiber.StatusBadRequest, "failed to get file from form")
	}

	err = h.ImageService.Create(c.Context(), fileHeader)
	if err != nil {
		log.Printf("%s: %v", codepath, err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create image")
	}

	c.Status(fiber.StatusCreated)
	return nil
}
