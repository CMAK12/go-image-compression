package http

import (
	"fmt"
	"io"
	"log"

	"go-image-compression/internal/model"

	"github.com/gofiber/fiber/v2"
)

const codepath = "controller/v1/http/image.go"

func (h *Handler) getImage(c *fiber.Ctx) (io.Reader, int, error) {
	var filter model.ListImageFilter

	if err := c.QueryParser(&filter); err != nil {
		log.Printf("%s: %v", codepath, err)
		return nil, fiber.StatusBadRequest, fmt.Errorf("%s: %w", codepath, err)
	}
	if filter.CompressPercent == 0 || filter.ID == "" {
		log.Printf("%s: missing required params: compress_percent=%d, id=%s", codepath, filter.CompressPercent, filter.ID)
		return nil, fiber.StatusBadRequest, fmt.Errorf("compress_percent and id are required")
	}

	image, err := h.Services.ImageService.Get(c.Context(), filter)
	if err != nil {
		log.Printf("%s: %v", codepath, err)
		return nil, fiber.StatusInternalServerError, fmt.Errorf("failed to get image")
	}

	return image, fiber.StatusOK, nil
}

func (h *Handler) createImage(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Printf("%s: %v", codepath, err)
		return fiber.NewError(fiber.StatusBadRequest, "failed to get file from form")
	}

	err = h.Services.ImageService.Create(c.Context(), fileHeader)
	if err != nil {
		log.Printf("%s: %v", codepath, err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create image")
	}

	c.Status(fiber.StatusCreated)
	return nil
}
