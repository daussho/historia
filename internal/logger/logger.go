package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	logger *zerolog.Logger
)

func Get() *zerolog.Logger {
	if logger != nil {
		return logger
	}

	level := zerolog.InfoLevel
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		lvl, err := zerolog.ParseLevel(logLevel)
		if err == nil {
			level = lvl
		}
	}

	lgr := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()

	logger = &lgr

	return logger
}

func Info() *zerolog.Event {
	return Get().Info()
}

func Error() *zerolog.Event {
	return Get().Error()
}

func Fatal() *zerolog.Event {
	return Get().Fatal()
}

func Panic() *zerolog.Event {
	return Get().Panic()
}

func Debug() *zerolog.Event {
	return Get().Debug()
}

func Trace() *zerolog.Event {
	return Get().Trace()
}

func Warn() *zerolog.Event {
	return Get().Warn()
}
