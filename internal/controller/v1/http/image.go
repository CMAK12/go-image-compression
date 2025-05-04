package http

import "github.com/gofiber/fiber/v2"

func (h *Handler) listImage(c *fiber.Ctx) error {
	return c.SendString("List Image")
}

func (h *Handler) createImage(c *fiber.Ctx) error {
	return c.SendString("Create Image")
}
