package app

import (
	"os"
	"testing"
)

func TestConf_parseConfig(t *testing.T) {
	t.Run("envValuesNotProvided_returnsConfigWithDefaultValues", func(t *testing.T) {
		config := parseConfig()

		if config.Addr != defaultAddr {
			t.Errorf("Expected %s, got %s", defaultAddr, config.Addr)
		}
		if config.allowOverlocking != defaultAllowOverlocking {
			t.Errorf("Expected %T, got %T", defaultAllowOverlocking, config.allowOverlocking)
		}
		if config.pipelinesCaseSensitive != defaultPipelinesCaseSensitive {
			t.Errorf("Expected %T, got %T", defaultPipelinesCaseSensitive, config.pipelinesCaseSensitive)
		}
	})

	const (
		addrValue = ":9000"
	)
	t.Run("envValuesProvided_returnsConfigWithProvidedValues", func(t *testing.T) {
		os.Setenv(addrEnvKey, addrValue)
		os.Setenv(allowOverlockingKey, "true")
		os.Setenv(pipelinesCaseSensitiveKey, "false")
		config := parseConfig()

		if config.Addr != addrValue {
			t.Errorf("Expected address %s, got %s", addrValue, config.Addr)
		}
		if config.allowOverlocking != true {
			t.Errorf("Expected %T, got %T", true, config.allowOverlocking)
		}
		if config.pipelinesCaseSensitive != false {
			t.Errorf("Expected %T, got %T", false, config.pipelinesCaseSensitive)
		}
		os.Clearenv()
	})
}
