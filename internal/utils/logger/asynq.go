package logger

import (
	"github.com/rs/zerolog"
)

type AsynqLogger struct {
	logger *zerolog.Logger
}

func NewAsynqLogger(logger *zerolog.Logger) *AsynqLogger {
	log := logger.With().Str("eventClass", "asynq").Str("event", "asynq.event").Logger()
	return &AsynqLogger{
		logger: &log,
	}
}

func (l *AsynqLogger) Debug(args ...interface{}) {
	l.logger.Debug().Msgf("%v", args...)
}

func (l *AsynqLogger) Info(args ...interface{}) {
	l.logger.Info().Msgf("%v", args...)
}

func (l *AsynqLogger) Warn(args ...interface{}) {
	l.logger.Warn().Msgf("%v", args...)
}

func (l *AsynqLogger) Error(args ...interface{}) {
	l.logger.Error().Msgf("%v", args...)
}

func (l *AsynqLogger) Fatal(args ...interface{}) {
	l.logger.Fatal().Msgf("%v", args...)
}
