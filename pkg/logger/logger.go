package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// Init initializes the logger with custom settings.
func Init(logLevel string) {
	log = logrus.New()

	// Set log output to stdout
	log.Out = os.Stdout

	// Set the log format
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Parse the log level from the string
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		// Default to WARN level if the provided log level is invalid
		level = logrus.WarnLevel
	}

	// Set the default log level (WARN/ERROR)
	log.SetLevel(level)
}

// GetLogger returns the configured logger instance.
func GetLogger() *logrus.Logger {
	return log
}
