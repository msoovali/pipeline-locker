package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/msoovali/pipeline-locker/internal/app"
)

func main() {
	htmlEngine := html.New("./views", ".html")
	router := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		Views:        htmlEngine,
	})
	router.Use(logger.New())
	router.Use(favicon.New(favicon.Config{
		File: "./views/favicon.ico",
	}))
	app := app.New(router)
	app.Log.Error.Fatal(router.Listen(app.Config.Addr))
}
