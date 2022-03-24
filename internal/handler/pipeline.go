package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/msoovali/pipeline-locker/internal/domain"
	"github.com/msoovali/pipeline-locker/internal/service"
)

type pipelineHandlers struct {
	service service.PipelineService
}

func NewPipelineHandlers(service service.PipelineService) *pipelineHandlers {
	return &pipelineHandlers{
		service: service,
	}
}

func (h *pipelineHandlers) Lock(c *fiber.Ctx) error {
	r := new(domain.Pipeline)
	if err := c.BodyParser(r); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.service.Lock(createImmutablePipeline(*r)); err != nil {
		return c.Status(fiber.StatusConflict).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *pipelineHandlers) Unlock(c *fiber.Ctx) error {
	r := new(domain.PipelineIdentifier)
	if err := c.BodyParser(r); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.service.Unlock(createImmutablePipelineIdentifier(*r)); err != nil {
		return c.Status(fiber.StatusConflict).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *pipelineHandlers) GetStatus(c *fiber.Ctx) error {
	allowed, err := h.service.IsDeployAllowed(domain.PipelineIdentifier{
		Project:     c.Params("project"),
		Environment: c.Params("environment"),
	})
	if err != nil {
		return c.Status(fiber.StatusConflict).SendString(err.Error())
	}
	if !allowed {
		return c.Status(fiber.StatusLocked).SendString("PIPELINE_IS_LOCKED")
	}

	return c.SendString("OK")
}

func (h *pipelineHandlers) GetLockedPipelines(c *fiber.Ctx) error {
	return c.JSON(h.service.GetLockedPipelines())
}

func (h *pipelineHandlers) Index(c *fiber.Ctx) error {
	pipelines := h.service.GetLockedPipelines()
	return c.Render("index", fiber.Map{
		"pipelines": pipelines,
	}, "layouts/main")
}

func (h *pipelineHandlers) LockAndRedirect(c *fiber.Ctx) error {
	r := new(domain.Pipeline)
	err := c.BodyParser(r)
	if err == nil {
		err = h.service.Lock(createImmutablePipeline(*r))
	}
	if err == nil {
		return c.Redirect("/", fiber.StatusSeeOther)
	}
	pipelines := h.service.GetLockedPipelines()

	return c.Render("index", fiber.Map{
		"err":       err,
		"pipelines": pipelines,
		"formInput": r,
	}, "layouts/main")
}

func createImmutablePipelineIdentifier(p domain.PipelineIdentifier) domain.PipelineIdentifier {
	return domain.PipelineIdentifier{
		Project:     utils.ImmutableString(p.Project),
		Environment: utils.ImmutableString(p.Environment),
	}
}

func createImmutablePipeline(p domain.Pipeline) domain.Pipeline {
	return domain.Pipeline{
		PipelineIdentifier: createImmutablePipelineIdentifier(p.PipelineIdentifier),
		LockedBy:           utils.ImmutableString(p.LockedBy),
	}
}
