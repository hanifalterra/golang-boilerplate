package logger

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"golang-boilerplate/internal/pkg/utils/auth"
)

// AppLogger wraps a zerolog.Logger for structured logging with context.
type AppLogger struct {
	logger *zerolog.Logger
}

// Context keys for storing logger and start time.
type loggerCtxKey struct{}
type startTimeCtxKey struct{}

// NewAppLogger initializes a new AppLogger with a unique event ID and attaches it to the context.
// It also records the current time for measuring elapsed time.
func NewAppLogger(ctx context.Context, baseLogger *zerolog.Logger) (context.Context, *AppLogger) {
	now := time.Now()
	eventID := uuid.NewString()
	ctx = context.WithValue(ctx, startTimeCtxKey{}, now)

	logger := baseLogger.With().
		Str("eventID", eventID).
		Logger().
		Hook(TracingHook{})

	// Store the logger in the context.
	ctx = context.WithValue(ctx, loggerCtxKey{}, &AppLogger{logger: &logger})

	return ctx, &AppLogger{logger: &logger}
}

// NewAppLoggerEcho initializes a new AppLogger with an Echo context.
// It includes user credentials and request ID in the logger metadata.
func NewAppLoggerEcho(c echo.Context, baseLogger *zerolog.Logger) (context.Context, *AppLogger) {
	now := time.Now()
	eventID := getRequestID(c)
	ctx := context.WithValue(context.Background(), startTimeCtxKey{}, now)

	logger := baseLogger.With().
		Str("eventID", eventID).
		Str("credential", auth.GetUser(c).Username).
		Logger().
		Hook(TracingHook{})

	// Store the logger in the context.
	ctx = context.WithValue(ctx, loggerCtxKey{}, &AppLogger{logger: &logger})

	return ctx, &AppLogger{logger: &logger}
}

var nopLogger *AppLogger

// init initializes a no-op logger to use as a fallback when no logger is found in the context.
func init() {
	nop := zerolog.Nop()
	nopLogger = &AppLogger{logger: &nop}
}

// FromContext retrieves the AppLogger from the context.
// If no logger is found, it returns a no-op logger that does nothing when called.
func FromContext(ctx context.Context) *AppLogger {
	if logger, ok := ctx.Value(loggerCtxKey{}).(*AppLogger); ok {
		return logger
	}
	return nopLogger
}

// getRequestID retrieves the request ID from Echo's request headers.
// If the request ID is not found in the request headers, it falls back to the response headers.
func getRequestID(c echo.Context) string {
	if reqID := c.Request().Header.Get(echo.HeaderXRequestID); reqID != "" {
		return reqID
	}
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// Debug logs a message at the Debug level with event class and event name metadata.
func (l *AppLogger) Debug(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.DebugLevel, eventClass, eventName, message, args...)
}

// Info logs a message at the Info level with event class and event name metadata.
func (l *AppLogger) Info(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.InfoLevel, eventClass, eventName, message, args...)
}

// Warn logs a message at the Warn level with event class and event name metadata.
func (l *AppLogger) Warn(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.WarnLevel, eventClass, eventName, message, args...)
}

// Error logs a message at the Error level with event class and event name metadata.
func (l *AppLogger) Error(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.ErrorLevel, eventClass, eventName, message, args...)
}

// Fatal logs a message at the Fatal level with event class and event name metadata,
// and terminates the application.
func (l *AppLogger) Fatal(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.FatalLevel, eventClass, eventName, message, args...)
}

// log handles the logging process at the specified level,
// adding event metadata such as event class and event name before logging the message.
func (l *AppLogger) log(ctx context.Context, level zerolog.Level, eventClass, eventName, message string, args ...interface{}) {
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

// TracingHook is a zerolog hook that adds tracing information (such as elapsed time in milliseconds) to log entries.
type TracingHook struct{}

// Run is invoked by zerolog for each log event. It adds the elapsed time since the logger was created
// if the start time is available in the context.
func (h TracingHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if ctx := e.GetCtx(); ctx != nil {
		if startTime, ok := ctx.Value(startTimeCtxKey{}).(time.Time); ok {
			elapsedTime := time.Since(startTime).Milliseconds()
			e.Int64("elapsedTime", elapsedTime)
		}
	}
}
