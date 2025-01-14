package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"golang-boilerplate/internal/pkg/config"
)

// New creates a new base logger and attaches system context.
func New(cfg *config.Config, serviceName string) *zerolog.Logger {
	// Map logger level from config
	var level zerolog.Level
	switch strings.ToLower(cfg.Logger.Level) {
	case "error":
		level = zerolog.ErrorLevel
	case "warn":
		level = zerolog.WarnLevel
	case "info":
		level = zerolog.InfoLevel
	case "debug":
		level = zerolog.DebugLevel
	default:
		level = zerolog.InfoLevel
	}

	// Set global logger level
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339Nano

	// Create and return the logger instance
	logger := zerolog.New(os.Stdout).
		Level(level).
		With().
		Str("systemName", cfg.App.Name).
		Str("systemVersion", cfg.App.Version).
		Str("serviceName", serviceName).
		Timestamp().
		Logger()

	return &logger
}
