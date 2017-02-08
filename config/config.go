package config

import (
	"os"
)

type AppConfig struct {
	MetricsPath   string
	MetricsPort   string
	ListenPort    string
	LogLevel      string
	BlueprintPath string
}

func Init() AppConfig {
	appConfig := AppConfig{
		MetricsPath:   GetEnv("METRICS_PATH", "/metrics"),
		MetricsPort:   GetEnv("METRICS_PORT", ":8090"),
		ListenPort:    GetEnv("LISTEN_PORT", "8080"),
		LogLevel:      GetEnv("LOG_LEVEL", "debug"),
		BlueprintPath: GetEnv("BLUEPRINT_PATH", "/blueprints/"),
	}

	return appConfig
}

// getEnv - Allows us to supply a fallback option if nothing specified
func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
