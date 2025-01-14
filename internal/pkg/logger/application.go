package logger

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"golang-boilerplate/internal/pkg/utils/auth"
)

// AppLogger defines the interface for structured logging with multiple log levels.
type AppLogger interface {
	Debug(ctx context.Context, eventClass, eventName, message string, args ...interface{})
	Info(ctx context.Context, eventClass, eventName, message string, args ...interface{})
	Warn(ctx context.Context, eventClass, eventName, message string, args ...interface{})
	Error(ctx context.Context, eventClass, eventName, message string, args ...interface{})
	Fatal(ctx context.Context, eventClass, eventName, message string, args ...interface{})
}

// appLogger is the concrete implementation of AppLogger using zerolog.
type appLogger struct {
	logger *zerolog.Logger
}

// startTimeCtxKey is a type-safe key for storing start time in context.
type startTimeCtxKey struct{}

// NewAppLogger initializes a new AppLogger with a unique event ID and attaches it to the context.
func NewAppLogger(ctx context.Context, baseLogger *zerolog.Logger) (context.Context, AppLogger) {
	now := time.Now()
	eventID := uuid.NewString()
	ctx = context.WithValue(ctx, startTimeCtxKey{}, now)

	logger := baseLogger.With().Str("eventID", eventID).Logger().Hook(TracingHook{})

	return ctx, &appLogger{logger: &logger}
}

// NewAppLoggerEcho initializes a new AppLogger with Echo context, including user credentials and request ID.
func NewAppLoggerEcho(c echo.Context, baseLogger *zerolog.Logger) (context.Context, AppLogger) {
	now := time.Now()
	eventID := getRequestID(c)
	ctx := context.WithValue(context.Background(), startTimeCtxKey{}, now)

	logger := baseLogger.With().
		Str("eventID", eventID).
		Str("credential", auth.GetUser(c).Username).
		Logger().
		Hook(TracingHook{})

	return ctx, &appLogger{logger: &logger}
}

// getRequestID retrieves the request ID from Echo's request or response headers.
// If the request ID is not found in the request headers, it falls back to the response headers.
func getRequestID(c echo.Context) string {
	if reqID := c.Request().Header.Get(echo.HeaderXRequestID); reqID != "" {
		return reqID
	}
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// Debug logs a message at the Debug level with context information.
func (l *appLogger) Debug(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.DebugLevel, eventClass, eventName, message, args...)
}

// Info logs a message at the Info level with context information.
func (l *appLogger) Info(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.InfoLevel, eventClass, eventName, message, args...)
}

// Warn logs a message at the Warn level with context information.
func (l *appLogger) Warn(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.WarnLevel, eventClass, eventName, message, args...)
}

// Error logs a message at the Error level with context information.
func (l *appLogger) Error(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.ErrorLevel, eventClass, eventName, message, args...)
}

// Fatal logs a message at the Fatal level, logs context information, and exits the application.
func (l *appLogger) Fatal(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.FatalLevel, eventClass, eventName, message, args...)
}

// log handles the logging process, including adding event metadata and formatting the message.
func (l *appLogger) log(ctx context.Context, level zerolog.Level, eventClass, eventName, message string, args ...interface{}) {
	logger := l.logger.With().
		Str("eventClass", eventClass).
		Str("event", eventName).
		Logger()

	switch level {
	case zerolog.DebugLevel:
		logger.Debug().Ctx(ctx).Msgf(message, args...)
	case zerolog.InfoLevel:
		logger.Info().Ctx(ctx).Msgf(message, args...)
	case zerolog.WarnLevel:
		logger.Warn().Ctx(ctx).Msgf(message, args...)
	case zerolog.ErrorLevel:
		logger.Error().Ctx(ctx).Msgf(message, args...)
	case zerolog.FatalLevel:
		logger.Fatal().Ctx(ctx).Msgf(message, args...)
	}
}

// TracingHook is a zerolog hook that adds tracing information (e.g., elapsed time) to log entries.
type TracingHook struct{}

// Run is invoked by zerolog for each log event and adds elapsed time (in milliseconds) if start time is available in context.
func (h TracingHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if ctx := e.GetCtx(); ctx != nil {
		if startTime, ok := ctx.Value(startTimeCtxKey{}).(time.Time); ok {
			elapsedTime := time.Since(startTime).Milliseconds()
			e.Int64("elapsedTime", elapsedTime)
		}
	}
}
