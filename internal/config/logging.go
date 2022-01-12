package config

import (
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/term"
)

// SetFormatter sets the appropriate log formatter for our terminal (JSON if we're in Docker):
func (l *LoggingConfig) SetFormatter(logger *logrus.Logger) {

	// Detect if we're running in a terminal:
	if term.IsTerminal(int(os.Stdout.Fd())) {
		logger.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	// Set the logging level specified in the config:
	loggingLevel, err := logrus.ParseLevel(l.Level)
	if err != nil {
		logger.WithError(err).Warn("Invalid log level")
		return
	}
	logger.SetLevel(loggingLevel)
}
