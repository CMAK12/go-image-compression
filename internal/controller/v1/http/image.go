package http

import "github.com/gofiber/fiber/v2"

func (h *Handler) ListImage(c *fiber.Ctx) error {
	return c.SendString("List Image")
}

func (h *Handler) CreateImage(c *fiber.Ctx) error {
	return c.SendString("Create Image")
}
