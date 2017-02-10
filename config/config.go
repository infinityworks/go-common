package config

import "os"

type AppConfig interface {
	MetricsPath() string
	MetricsPort() string
	ListenPort() string
	LogLevel() string
}

type BaseConfig struct {
	metricsPath string
	metricsPort string
	listenPort  string
	logLevel    string
}

func (c BaseConfig) MetricsPath() string {
	return c.metricsPath
}

func (c BaseConfig) MetricsPort() string {
	return c.metricsPort
}

func (c BaseConfig) ListenPort() string {
	return c.listenPort
}

func (c BaseConfig) LogLevel() string {
	return c.logLevel
}

func Init() BaseConfig {

	appConfig := BaseConfig{
		metricsPath: GetEnv("METRICS_PATH", "/metrics"),
		metricsPort: GetEnv("METRICS_PORT", ":8090"),
		listenPort:  GetEnv("LISTEN_PORT", "8080"),
		logLevel:    GetEnv("LOG_LEVEL", "debug"),
	}
	return appConfig
}

// GetEnv - Allows us to supply a fallback option if nothing specified
func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
