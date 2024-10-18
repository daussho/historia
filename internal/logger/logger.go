package logger

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var (
	logger *log.Logger
)

func Log() *log.Logger {
	if logger != nil {
		return logger
	}

	level, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = log.InfoLevel
	}

	logger = log.NewWithOptions(os.Stdout, log.Options{
		TimeFormat:      time.RFC3339,
		Level:           level,
		ReportTimestamp: true,
		ReportCaller:    true,
	})

	return logger
}
