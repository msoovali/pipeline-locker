package app

import (
	"log"
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

func parseConfig() *ApplicationConfig {
	addr := getEnv(addrEnvKey, defaultAddr)
	allowOverlocking := getEnvBool(allowOverlockingKey, defaultAllowOverlocking)
	pipelinesCaseSensitive := getEnvBool(pipelinesCaseSensitiveKey, defaultPipelinesCaseSensitive)

	return &ApplicationConfig{
		Addr:                   addr,
		allowOverlocking:       allowOverlocking,
		pipelinesCaseSensitive: pipelinesCaseSensitive,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if valueString, ok := os.LookupEnv(key); ok {
		value, err := strconv.Atoi(valueString)
		if err == nil {
			return value
		}
		log.Printf("Failed to convert %s env value %s to int. Falling back to default %d", key, valueString, fallback)
	}

	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if valueString, ok := os.LookupEnv(key); ok {
		value, err := strconv.ParseBool(valueString)
		if err == nil {
			return value
		}
		log.Printf("Failed to convert %s env value %s to boolean. Falling back to default %T", key, valueString, fallback)
	}

	return fallback
}
