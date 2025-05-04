package http

import (
	"io"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
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

func ResponseWrapper(handler func(c *fiber.Ctx) (io.Reader, int, error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		response, statusCode, err := handler(c)
		if err != nil {
			log.Printf("%s: %v", codepath, err)
			return writeError(c, statusCode, err)
		}
		defer closeIOReader(response)

		c.Set("Content-Type", "image/png")
		c.Status(statusCode)

		_, err = io.Copy(c.Response().BodyWriter(), response)
		if err != nil {
			log.Printf("%s: %v", codepath, err)
			return writeError(c, fiber.StatusInternalServerError, err)
		}

		return nil
	}
}

func writeError(c *fiber.Ctx, statusCode int, err error) error {
	c.Set("Content-Type", "application/json")
	return c.Status(statusCode).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func closeIOReader(reader io.Reader) {
	if closer, ok := reader.(io.Closer); ok {
		err := closer.Close()
		if err != nil {
			log.Printf("%s: %v", codepath, err)
		}
	}
}
