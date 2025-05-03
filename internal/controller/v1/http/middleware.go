package http

import (
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
