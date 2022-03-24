package handler

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func TestHealthHandler_HealthCheck(t *testing.T) {
	t.Run("responds_statusAndMessageOk", func(t *testing.T) {
		handler := NewHealthHandlers()
		app := fiber.New()
		c := app.AcquireCtx(&fasthttp.RequestCtx{})
		defer app.ReleaseCtx(c)
		err := handler.HealthCheck(c)
		if err != nil {
			t.Errorf("Expected error nil, got %v", err)
		}
		if c.Response().StatusCode() != fiber.StatusOK {
			t.Errorf("Expected status code %d, got %d", fiber.StatusOK, c.Response().StatusCode())
		}
		if string(c.Response().Body()) != "OK" {
			t.Errorf("Expected 'OK' response, but got '%s'", string(c.Response().Body()))
		}
	})
}
