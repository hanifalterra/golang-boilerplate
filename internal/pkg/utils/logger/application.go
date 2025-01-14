package logger

import (
	"context"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"golang-boilerplate/internal/pkg/utils/auth"
)

// Interface defines the logging interface.
type LoggerInterface interface {
	Debug(ctx context.Context, eventClass, eventName, message string, args ...interface{})
	Info(ctx context.Context, eventClass, eventName, message string, args ...interface{})
	Warn(ctx context.Context, eventClass, eventName, message string, args ...interface{})
	Error(ctx context.Context, eventClass, eventName, message string, args ...interface{})
	Fatal(ctx context.Context, eventClass, eventName, message string, args ...interface{})
}

// Logger represents the logger structure.
type Logger struct {
	logger *zerolog.Logger
}

// Interface checker
var (
	_ LoggerInterface = (*Logger)(nil)
)

// contextKey is a private type for context keys to avoid collisions.
type contextKey string

const (
	loggerCtxKey    contextKey = "logger"
	eventIDCtxKey   contextKey = "eventID"
	userCtxKey      contextKey = "user"
	startTimeCtxKey contextKey = "startTime"
)

// NewContextWithLogger creates a new logger and attaches it to the context.
func NewContextWithLogger(ctx context.Context, log *zerolog.Logger, eventID string) (context.Context, *Logger) {
	if eventID == "" {
		eventID = uuid.NewString()
	}

	ctx = context.WithValue(ctx, eventIDCtxKey, eventID)
	ctx = context.WithValue(ctx, startTimeCtxKey, time.Now())

	logger := log.With().Str("eventID", eventID).Logger().Hook(TracingHook{})
	l := &Logger{logger: &logger}

	return context.WithValue(ctx, loggerCtxKey, l), l
}

// NewContextEchoWithLogger creates a logger with additional info from Echo context.
func NewContextEchoWithLogger(c echo.Context, log *zerolog.Logger) (context.Context, *Logger) {
	eventID := getRequestID(c)
	user := auth.GetUser(c)
	now := time.Now()

	ctx := context.WithValue(context.Background(), eventIDCtxKey, eventID)
	ctx = context.WithValue(ctx, userCtxKey, user)
	ctx = context.WithValue(ctx, startTimeCtxKey, now)

	logger := log.With().
		Str("eventID", eventID).
		Str("credential", user.Username).
		Logger().
		Hook(TracingHook{})

	l := &Logger{logger: &logger}
	return context.WithValue(ctx, loggerCtxKey, l), l
}

// Debug logs a debug message.
func (l *Logger) Debug(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.DebugLevel, eventClass, eventName, message, args...)
}

// Info logs an info message.
func (l *Logger) Info(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.InfoLevel, eventClass, eventName, message, args...)
}

// Warn logs a warning message.
func (l *Logger) Warn(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.WarnLevel, eventClass, eventName, message, args...)
}

// Error logs an error message.
func (l *Logger) Error(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.ErrorLevel, eventClass, eventName, message, args...)
}

// Fatal logs a fatal message and exits the application.
func (l *Logger) Fatal(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, zerolog.FatalLevel, eventClass, eventName, message, args...)
	os.Exit(1)
}

// log is the generic log handler for all levels.
func (l *Logger) log(ctx context.Context, level zerolog.Level, eventClass, eventName, message string, args ...interface{}) {
	event := l.logger.WithLevel(level).Str("eventClass", eventClass).Str("event", eventName)

	// Add elapsed time if available
	if startTime, ok := ctx.Value(startTimeCtxKey).(time.Time); ok {
		elapsedTime := time.Since(startTime).Milliseconds()
		event = event.Int64("elapsedTime", elapsedTime)
	}

	// Add user information if available
	if user, ok := ctx.Value(userCtxKey).(auth.User); ok {
		event = event.Str("user", user.Username)
	}

	if len(args) == 0 {
		event.Ctx(ctx).Msg(message)
	} else {
		event.Ctx(ctx).Msgf(message, args...)
	}
}

// TracingHook adds elapsed time information to log entries.
type TracingHook struct{}

// Run adds the elapsed time to the log event.
func (h TracingHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	ctx := e.GetCtx()
	if startTime, ok := ctx.Value(startTimeCtxKey).(time.Time); ok {
		elapsedTime := time.Since(startTime).Milliseconds()
		e.Int64("elapsedTime", elapsedTime)
	}
}
