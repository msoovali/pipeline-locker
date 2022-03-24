package handler

import "github.com/gofiber/fiber/v2"

type PipelineHandlers interface {
	Lock(c *fiber.Ctx) error
	Unlock(c *fiber.Ctx) error
	GetStatus(c *fiber.Ctx) error
	GetLockedPipelines(c *fiber.Ctx) error
	Index(c *fiber.Ctx) error
	LockAndRedirect(c *fiber.Ctx) error
}

type HealthHandlers interface {
	HealthCheck(c *fiber.Ctx) error
}
