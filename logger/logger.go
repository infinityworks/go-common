package logger

import (
	log "github.com/Sirupsen/logrus"
	"github.com/infinityworksltd/go-common/config"
)

func Start(config config.AppConfig) (logger *log.Logger) {

	logger = log.New()
	logger.Formatter = &log.JSONFormatter{}
	setLogLevel(config.LogLevel(), logger)

	return logger
}

// setLogLevel - Sets the log level based on the passed argument.
func setLogLevel(level string, l *log.Logger) {
	switch level {
	case "debug":
		l.Level = log.DebugLevel
		break
	case "info":
		l.Level = log.InfoLevel
		break
	case "warn":
		l.Level = log.WarnLevel
		break
	case "fatal":
		l.Level = log.FatalLevel
		break
	case "panic":
		l.Level = log.PanicLevel
		break
	}
}
