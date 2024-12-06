package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"golang-boilerplate/internal/pkg/config"
)

// Creates New Base logger and attach system context.
func New(cfg *config.Config, serviceName string) *zerolog.Logger {
	var l zerolog.Level

	switch strings.ToLower(cfg.Logger.Level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	// Remark for when need to shorten filename
	// zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
	// 	_, b, _, _ := runtime.Caller(0)
	// 	projectRoot := filepath.Dir(b)
	// 	projectRoot = strings.Replace(projectRoot, "/utils/logger", "", 1)
	// 	file = strings.Replace(file, projectRoot, "", 1)
	// 	return file + ":" + strconv.Itoa(line)
	// }

	zerolog.SetGlobalLevel(l)
	zerolog.TimeFieldFormat = time.RFC3339Nano
	logger := zerolog.New(os.Stdout).Level(l).With().
		Str("systemName", cfg.App.Name).
		Str("systemVersion", cfg.App.Version).
		Str("serviceName", serviceName).
		Timestamp().
		Logger()

	return &logger
}
