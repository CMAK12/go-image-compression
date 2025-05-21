package http

import (
	"fmt"
	"log"
	"mime/multipart"

	"go-image-compression/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/nordew/go-errx"
)

const codepath = "controller/v1/http/image.go"

func (h *Handler) getImage(c *fiber.Ctx) (multipart.File, error) {
	var filter model.ListImageFilter

	if err := c.QueryParser(&filter); err != nil {
		return nil, errx.NewBadRequest().WithDescription("failed to parse query params")
	}
	if filter.ID == "" {
		log.Printf("%s: missing required params: id=%s", codepath, filter.ID)
		return nil, errx.NewBadRequest().WithDescription("id field is required")
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
