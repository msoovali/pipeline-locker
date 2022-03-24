package app

import (
	"github.com/gofiber/fiber/v2"
)

func (a *Application) registerRoutes(router *fiber.App) {
	router.Get("/health", a.Handlers.HealthHandlers.HealthCheck)
	router.Get("/", a.Handlers.PipelineHandlers.Index)
	router.Post("/", a.Handlers.PipelineHandlers.LockAndRedirect)

	v1 := router.Group("/v1")
	{
		v1.Post("/pipeline/lock", a.Handlers.PipelineHandlers.Lock)
		v1.Put("/pipeline/unlock", a.Handlers.PipelineHandlers.Unlock)
		v1.Get("/pipeline/status/project/:project/environment/:environment", a.Handlers.PipelineHandlers.GetStatus)
		v1.Get("/pipelines/locked", a.Handlers.PipelineHandlers.GetLockedPipelines)
	}
}
