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
	redisVersionKey               = "REDIS_VERSION"
	redisAddr                     = "REDIS_ADDR"
	defaultRedisAddr              = "localhost:6379"
	redisUsername                 = "REDIS_USERNAME"
	defaultRedisUsername          = ""
	redisPassword                 = "REDIS_PASSWORD"
	defaultRedisPassword          = ""
)

type ApplicationConfig struct {
	Addr                   string
	allowOverlocking       bool
	pipelinesCaseSensitive bool
	redisConfig            *redisConfig
}

type redisConfig struct {
	version  int
	addr     string
	username string
	password string
}

func (a *Application) parseConfig() {
	a.Config = &ApplicationConfig{
		Addr:                   a.getEnv(addrEnvKey, defaultAddr),
		allowOverlocking:       a.getEnvBool(allowOverlockingKey, defaultAllowOverlocking),
		pipelinesCaseSensitive: a.getEnvBool(pipelinesCaseSensitiveKey, defaultPipelinesCaseSensitive),
	}

	redisVersion := a.getEnvInt(redisVersionKey, 0)
	if redisVersion != 0 {
		a.parseRedisConfig(redisVersion)
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

func (a *Application) parseRedisConfig(version int) {
	if version != 6 && version != 7 {
		a.Log.Error.Printf("Redis version %d is not supported, falling back to memory based repository. Redis versions 6 and 7 are supported!", version)
		return
	}

	a.Config.redisConfig = &redisConfig{
		version:  version,
		addr:     a.getEnv(redisAddr, defaultRedisAddr),
		username: a.getEnv(redisUsername, defaultRedisUsername),
		password: a.getEnv(redisPassword, defaultRedisPassword),
	}
}
