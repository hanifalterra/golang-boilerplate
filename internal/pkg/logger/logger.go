package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// New creates and returns a new base logger instance with system context.
// It initializes the logger level based on the provided configuration and attaches
// additional metadata such as system name, version, and service name.
func New(levelStr string, systemName string, systemVersion string, serviceName string) *zerolog.Logger {
	// Map the logger level from the configuration.
	var level zerolog.Level
	switch strings.ToLower(levelStr) {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	default:
		level = zerolog.InfoLevel // Default to Info level if not specified.
	}

	// Set the global logger level for zerolog.
	zerolog.SetGlobalLevel(level)

	// Set the time field format for log timestamps to RFC 3339 with nanosecond precision.
	zerolog.TimeFieldFormat = time.RFC3339Nano

	// Create the logger instance with the specified level and additional context.
	log.Logger = zerolog.New(os.Stdout).
		With().
		Str("systemName", systemName).
		Str("systemVersion", systemVersion).
		Str("serviceName", serviceName).
		Timestamp(). // Automatically include a timestamp in each log entry.
		Logger()

	return &log.Logger
}
