package app

import (
	redis_pkg_v8 "github.com/go-redis/redis/v8"
	redis_pkg_v9 "github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/msoovali/pipeline-locker/internal/domain"
	"github.com/msoovali/pipeline-locker/internal/handler"
	"github.com/msoovali/pipeline-locker/internal/logger"
	"github.com/msoovali/pipeline-locker/internal/repository/memory"
	redis_v6 "github.com/msoovali/pipeline-locker/internal/repository/redis/v6"
	redis_v7 "github.com/msoovali/pipeline-locker/internal/repository/redis/v7"
	"github.com/msoovali/pipeline-locker/internal/service"
)

type repositories struct {
	PipelineRepository domain.PipelineRepository
}

type services struct {
	PipelineService domain.PipelineService
}

type handlers struct {
	HealthHandlers   handler.HealthHandlers
	PipelineHandlers handler.PipelineHandlers
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
	var repositories *repositories
	redisConfig := a.Config.redisConfig
	if redisConfig != nil {
		if redisConfig.version == 6 {
			repositories = initRedis6Repositories(a.Config)
		} else if redisConfig.version == 7 {
			repositories = initRedis7Repositories(a.Config)
		}
	}
	if repositories == nil {
		repositories = initInMemoryRepositories(a.Config)
	}
	a.Repositories = repositories
}

func initInMemoryRepositories(config *ApplicationConfig) *repositories {
	return &repositories{
		PipelineRepository: memory.NewPipelineRepository(config.pipelinesCaseSensitive),
	}
}

func initRedis6Repositories(config *ApplicationConfig) *repositories {
	client := initRedis6Client(config.redisConfig)
	return &repositories{
		PipelineRepository: redis_v6.NewPipelineRepository(client, config.pipelinesCaseSensitive),
	}
}

func initRedis7Repositories(config *ApplicationConfig) *repositories {
	client := initRedis7Client(config.redisConfig)
	return &repositories{
		PipelineRepository: redis_v7.NewPipelineRepository(client, config.pipelinesCaseSensitive),
	}
}

func initRedis6Client(config *redisConfig) *redis_pkg_v8.Client {
	return redis_pkg_v8.NewClient(&redis_pkg_v8.Options{
		Addr:     config.addr,
		Username: config.username,
		Password: config.password,
	})
}

func initRedis7Client(config *redisConfig) *redis_pkg_v9.Client {
	return redis_pkg_v9.NewClient(&redis_pkg_v9.Options{
		Addr:     config.addr,
		Username: config.username,
		Password: config.password,
	})
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
