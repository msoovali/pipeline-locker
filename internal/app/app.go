package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/msoovali/pipeline-locker/internal/handler"
	"github.com/msoovali/pipeline-locker/internal/logger"
	"github.com/msoovali/pipeline-locker/internal/repository"
	"github.com/msoovali/pipeline-locker/internal/repository/memory"
	"github.com/msoovali/pipeline-locker/internal/service"
)

type repositories struct {
	PipelineRepository repository.PipelineRepository
}

type services struct {
	PipelineService service.PipelineService
}

type handlers struct {
	HealthHandlers   handler.HealthHandlers
	PipelineHandlers handler.PipelineHandlers
}

type IApplication interface {
	initRepositories()
	initServices()
	initHandlers()
	registerRoutes(router *fiber.App)
}

type Application struct {
	Log          *logger.Logger
	Config       *ApplicationConfig
	Repositories *repositories
	Services     *services
	Handlers     *handlers
}

func New(router *fiber.App) *Application {
	app := &Application{
		Log: logger.New(),
	}
	app.parseConfig()
	app.initRepositories()
	app.initServices()
	app.initHandlers()
	app.registerRoutes(router)

	return app
}

func (a *Application) initRepositories() {
	a.Repositories = &repositories{
		PipelineRepository: memory.NewPipelineRepository(a.Config.pipelinesCaseSensitive),
	}
}

func (a *Application) initServices() {
	a.Services = &services{
		PipelineService: service.NewPipelineService(a.Repositories.PipelineRepository, a.Config.allowOverlocking),
	}
}

func (a *Application) initHandlers() {
	a.Handlers = &handlers{
		HealthHandlers:   handler.NewHealthHandlers(),
		PipelineHandlers: handler.NewPipelineHandlers(a.Services.PipelineService),
	}
}
