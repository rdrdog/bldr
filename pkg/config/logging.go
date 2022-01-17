package config

import (
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/term"
)

// Configure sets the appropriate log formatter for our terminal (JSON if we're in Docker), and sets up masking and log levels
func (l *LoggingConfig) Configure(logger *logrus.Logger) {
	l.Masker = &MaskingHook{MaskedValue: "***"}

	// Detect if we're running in a terminal:
	if term.IsTerminal(int(os.Stdout.Fd())) {
		logger.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	logger.AddHook(l.Masker)

	// Set the logging level specified in the config:
	loggingLevel, err := logrus.ParseLevel(l.Level)
	if err != nil {
		logger.WithError(err).Warn("Invalid log level")
		return
	}
	logger.SetLevel(loggingLevel)
}

// Adds the specified secrets to the logging secret mask so that it's not emitted in the output
func (l *LoggingConfig) AddToSecretMask(secret string) {
	l.Masker.AddToMaskList(secret)
}
