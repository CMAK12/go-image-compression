package http

import (
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nordew/go-errx"
)

func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		log.Printf("Started: %s %s", c.Method(), c.Path())

		err := c.Next()

		log.Printf("Completed %s in %v", c.Path(), time.Since(start))
		return err
	}
}

func ResponseWrapper(handler func(c *fiber.Ctx) (multipart.File, error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		response, err := handler(c)
		if err != nil {
			return handleError(c, err)
		}

		return displayImage(c, response)
	}
}

func handleError(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	log.Printf("error occurred: %s: %v", codepath, err)

	switch {
	case errx.IsCode(err, errx.NotFound):
		return writeError(c, fiber.StatusNotFound, err)
	case errx.IsCode(err, errx.BadRequest):
		return writeError(c, fiber.StatusBadRequest, err)
	case errx.IsCode(err, errx.Internal):
		return writeError(c, fiber.StatusInternalServerError, err)
	case errx.IsCode(err, errx.Unauthorized):
		return writeError(c, fiber.StatusUnauthorized, err)
	case errx.IsCode(err, errx.Forbidden):
		return writeError(c, fiber.StatusForbidden, err)
	default:
		return writeError(c, fiber.StatusInternalServerError, fmt.Errorf("unexpected error: %v", err))
	}
}

func writeError(c *fiber.Ctx, statusCode int, err error) error {
	response := fiber.Map{
		"success": false,
		"error":   err.Error(),
	}

	return c.Status(statusCode).JSON(response)
}

func displayImage(c *fiber.Ctx, image multipart.File) error {
	c.Set("Content-Type", "image/jpeg")
	c.Set("Content-Disposition", "inline")
	return c.SendStream(image)
}
