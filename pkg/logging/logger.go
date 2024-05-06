package logging

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"sync"
)

const (
	EnvLogMode  = "LOG_MODE"
	EnvLogLevel = "LOG_LEVEL"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

var (
	defaultLogger     *slog.Logger
	defaultLoggerOnce sync.Once
)

func convertLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func FromEnv() *slog.Logger {
	develop := strings.ToLower(strings.TrimSpace(os.Getenv(EnvLogMode))) == "develop"
	return New(develop, os.Getenv(EnvLogLevel))
}

func New(develop bool, level string) *slog.Logger {
	var handler slog.Handler
	if develop {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     convertLevel(level),
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: convertLevel(level),
		})
	}
	return slog.New(handler)
}

func DefaultLogger() *slog.Logger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = FromEnv()
	})
	return defaultLogger
}

type contextKey string

const loggerKey contextKey = "logger"

func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return DefaultLogger()
}
