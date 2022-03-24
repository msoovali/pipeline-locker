package handler

import "github.com/gofiber/fiber/v2"

type healthHandlers struct {
}

func NewHealthHandlers() *healthHandlers {
	return &healthHandlers{}
}

func (h *healthHandlers) HealthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}
