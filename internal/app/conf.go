package app

import (
	"os"
	"strconv"
)

const (
	addrEnvKey                    = "ADDR"
	defaultAddr                   = ":8080"
	allowOverlockingKey           = "ALLOW_OVERLOCKING"
	defaultAllowOverlocking       = false
	pipelinesCaseSensitiveKey     = "PIPELINES_CASE_SENSITIVE"
	defaultPipelinesCaseSensitive = true
)

type ApplicationConfig struct {
	Addr                   string
	allowOverlocking       bool
	pipelinesCaseSensitive bool
}

func (a *Application) parseConfig() {
	addr := a.getEnv(addrEnvKey, defaultAddr)
	allowOverlocking := a.getEnvBool(allowOverlockingKey, defaultAllowOverlocking)
	pipelinesCaseSensitive := a.getEnvBool(pipelinesCaseSensitiveKey, defaultPipelinesCaseSensitive)

	a.Config = &ApplicationConfig{
		Addr:                   addr,
		allowOverlocking:       allowOverlocking,
		pipelinesCaseSensitive: pipelinesCaseSensitive,
	}
}

func (a *Application) getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func (a *Application) getEnvInt(key string, fallback int) int {
	if valueString, ok := os.LookupEnv(key); ok {
		value, err := strconv.Atoi(valueString)
		if err == nil {
			return value
		}
		a.Log.Error.Printf("Failed to convert %s env value %s to int. Falling back to default %d", key, valueString, fallback)
	}

	return fallback
}

func (a *Application) getEnvBool(key string, fallback bool) bool {
	if valueString, ok := os.LookupEnv(key); ok {
		value, err := strconv.ParseBool(valueString)
		if err == nil {
			return value
		}
		a.Log.Error.Printf("Failed to convert %s env value %s to boolean. Falling back to default %T", key, valueString, fallback)
	}

	return fallback
}
