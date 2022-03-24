package app

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestConf_parseConfig(t *testing.T) {
	t.Run("envValuesNotProvided_returnsConfigWithDefaultValues", func(t *testing.T) {
		app := New(fiber.New())
		app.parseConfig()

		if app.Config.Addr != defaultAddr {
			t.Errorf("Expected %s, got %s", defaultAddr, app.Config.Addr)
		}
		if app.Config.allowOverlocking != defaultAllowOverlocking {
			t.Errorf("Expected %T, got %T", defaultAllowOverlocking, app.Config.allowOverlocking)
		}
		if app.Config.pipelinesCaseSensitive != defaultPipelinesCaseSensitive {
			t.Errorf("Expected %T, got %T", defaultPipelinesCaseSensitive, app.Config.pipelinesCaseSensitive)
		}
	})

	const (
		addrValue = ":9000"
	)
	t.Run("envValuesProvided_returnsConfigWithProvidedValues", func(t *testing.T) {
		os.Setenv(addrEnvKey, addrValue)
		os.Setenv(allowOverlockingKey, "true")
		os.Setenv(pipelinesCaseSensitiveKey, "false")
		app := New(fiber.New())
		app.parseConfig()

		if app.Config.Addr != addrValue {
			t.Errorf("Expected address %s, got %s", addrValue, app.Config.Addr)
		}
		if app.Config.allowOverlocking != true {
			t.Errorf("Expected %T, got %T", true, app.Config.allowOverlocking)
		}
		if app.Config.pipelinesCaseSensitive != false {
			t.Errorf("Expected %T, got %T", false, app.Config.pipelinesCaseSensitive)
		}
		os.Clearenv()
	})
}
