package logger

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// LoggerInterface defines the contract for the logger.
type AppLogger interface {
	Debug(ctx context.Context, eventClass, eventName, message string, args ...interface{})
	Info(ctx context.Context, eventClass, eventName, message string, args ...interface{})
	Warn(ctx context.Context, eventClass, eventName, message string, args ...interface{})
	Error(ctx context.Context, eventClass, eventName, message string, args ...interface{})
	Fatal(ctx context.Context, eventClass, eventName, message string, args ...interface{})
}

// Logger encapsulates the zerolog logger.
type appLogger struct {
	logger *zerolog.Logger
}

type startTimeCtxKey struct{}

// NewContextWithLogger creates a new context with an attached logger instance.
// It generates a unique event ID if none is provided.
func NewAppLogger(ctx context.Context, l *zerolog.Logger, eventID string) (context.Context, AppLogger) {
	if eventID == "" {
		eventID = uuid.NewString()
	}

	ctx = context.WithValue(ctx, startTimeCtxKey{}, time.Now())

	logger := l.With().Str("eventID", eventID).Logger().Hook(TracingHook{})

	return ctx, &appLogger{logger: &logger}
}

// NewContextEchoWithLogger creates a context with a logger using Echo context values.
// It extracts request ID and user credentials from the Echo context.
func NewAppLoggerEcho(c echo.Context, l *zerolog.Logger, username string) (context.Context, AppLogger) {
	now := time.Now()
	eventID := getRequestID(c)

	ctx := context.WithValue(context.Background(), startTimeCtxKey{}, now)

	logger := l.With().
		Str("eventID", eventID).
		Str("credential", username).
		Logger().
		Hook(TracingHook{})

	return ctx, &appLogger{logger: &logger}
}

func getRequestID(c echo.Context) string {
	// Fetch X-Request-ID from request header, fallback to response header if not present
	if reqID := c.Request().Header.Get(echo.HeaderXRequestID); reqID != "" {
		return reqID
	}
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// Debug logs a message at Debug level.
func (l *appLogger) Debug(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.DebugLevel, eventClass, eventName, message, args...)
}

// Info logs a message at Info level.
func (l *appLogger) Info(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.InfoLevel, eventClass, eventName, message, args...)
}

// Warn logs a message at Warn level.
func (l *appLogger) Warn(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.WarnLevel, eventClass, eventName, message, args...)
}

// Error logs a message at Error level.
func (l *appLogger) Error(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.ErrorLevel, eventClass, eventName, message, args...)
}

// Fatal logs a message at Fatal level and immediately exits the application.
func (l *appLogger) Fatal(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.FatalLevel, eventClass, eventName, message, args...)
}

// log handles the common logic for logging at different levels.
func (l *appLogger) log(ctx context.Context, level zerolog.Level, eventClass, eventName, message string, args ...interface{}) {
	var event *zerolog.Event
	switch level {
	case zerolog.DebugLevel:
		event = l.logger.Debug().Str("eventClass", eventClass).Str("event", eventName)
	case zerolog.InfoLevel:
		event = l.logger.Info().Str("eventClass", eventClass).Str("event", eventName)
	case zerolog.WarnLevel:
		event = l.logger.Warn().Str("eventClass", eventClass).Str("event", eventName)
	case zerolog.ErrorLevel:
		event = l.logger.Error().Str("eventClass", eventClass).Str("event", eventName)
	case zerolog.FatalLevel:
		event = l.logger.Fatal().Str("eventClass", eventClass).Str("event", eventName)
	default:
		event = l.logger.Info().Str("eventClass", eventClass).Str("event", eventName)
	}

	if len(args) == 0 {
		event.Ctx(ctx).Msg(message)
	} else {
		event.Ctx(ctx).Msgf(message, args...)
	}
}

// TracingHook adds tracing information (e.g., elapsed time) to each log entry.
type TracingHook struct{}

// Run is executed by zerolog for each log event, adding elapsed time to the event.
func (h TracingHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	ctx := e.GetCtx()
	if startTime, ok := ctx.Value(startTimeCtxKey{}).(time.Time); ok {
		elapsedTime := time.Since(startTime).Milliseconds()
		e.Int64("elapsedTime", elapsedTime)
	}
}
