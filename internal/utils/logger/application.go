package logger

import (
	"context"
	"os"
	"time"

	auth "golang-boilerplate/internal/utils/auth"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// Interface -.
type Interface interface {
	Debug(ctx context.Context, eventClass, eventName, message string, args ...interface{})
	Info(ctx context.Context, eventClass, eventName, message string, args ...interface{})
	Warn(ctx context.Context, eventClass, eventName, message string, args ...interface{})
	Error(ctx context.Context, eventClass, eventName, message string, args ...interface{})
	Fatal(ctx context.Context, eventClass, eventName, message string, args ...interface{})
}

// Logger -.
type Logger struct {
	logger *zerolog.Logger
}

var _ Interface = (*Logger)(nil)

var disabledLogger *Logger

type loggerCtxKey struct{}
type ContextKeyEventID struct{}
type ContextKeyUser struct{}
type ContextKeyStartTime struct{}

func init() {
	disabledLogger = Nop()
}

// Disabled Logger for no logger in context fallback.
func Nop() *Logger {
	l := zerolog.New(nil).Level(zerolog.Disabled)
	return &Logger{
		logger: &l,
	}
}

// Get Logger from context, return Disabled logger when not found.
func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(loggerCtxKey{}).(*Logger); ok {
		return logger
	}

	return disabledLogger
}

func ContextResetTime(ctx context.Context) context.Context {
	return context.WithValue(ctx, ContextKeyStartTime{}, time.Now())
}

// Create new logger with context.
func NewContextWithLogger(ctx context.Context, log *zerolog.Logger, eventID string) (context.Context, *Logger) {
	if eventID == "" {
		uID, err := uuid.NewRandom()
		if err != nil {
			uID = uuid.New()
		}
		eventID = uID.String()
	}
	ctx = context.WithValue(ctx, ContextKeyEventID{}, eventID)
	ctx = context.WithValue(ctx, ContextKeyStartTime{}, time.Now())

	// skipFrameCount := 2
	logger := log.With().
		// CallerWithSkipFrameCount(
		// 	zerolog.CallerSkipFrameCount+skipFrameCount).
		Str("eventID", eventID).
		Logger()
	logger = logger.Hook(TracingHook{})
	l := &Logger{
		logger: &logger,
	}

	ctx = context.WithValue(ctx, loggerCtxKey{}, l)
	return ctx, l
}

// Create new logger with additional info in context.
func NewContextEchoWithLogger(c echo.Context, log *zerolog.Logger) (context.Context, *Logger) {
	eventID := getEchoRequestID(c)
	user := auth.GetUser(c)
	now := time.Now()
	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyEventID{}, eventID)
	ctx = context.WithValue(ctx, ContextKeyUser{}, user)
	ctx = context.WithValue(ctx, ContextKeyStartTime{}, now)
	username := user.Username

	// skipFrameCount := 2
	logger := log.With().
		// CallerWithSkipFrameCount(
		// 	zerolog.CallerSkipFrameCount+skipFrameCount).
		Str("eventID", eventID).
		Str("credential", username).
		Logger()
	logger = logger.Hook(TracingHook{})
	l := &Logger{
		logger: &logger,
	}
	ctx = context.WithValue(ctx, loggerCtxKey{}, l)
	return ctx, l
}

// Debug -.
func (l *Logger) Debug(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, "debug", eventClass, eventName, message, args...)
}

// Info -.
func (l *Logger) Info(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, "info", eventClass, eventName, message, args...)
}

// Warn -.
func (l *Logger) Warn(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, "warn", eventClass, eventName, message, args...)
}

// Error -.
func (l *Logger) Error(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, "error", eventClass, eventName, message, args...)
}

// Fatal -.
func (l *Logger) Fatal(ctx context.Context, eventClass, eventName, message string, args ...interface{}) {
	l.log(ctx, "fatal", eventClass, eventName, message, args...)

	os.Exit(1)
}

func (l *Logger) log(ctx context.Context, level, eventClass, eventName, message string, args ...interface{}) {
	var el *zerolog.Event
	switch level {
	case "debug":
		el = l.logger.Debug().Str("eventClass", eventClass).Str("event", eventName)
	case "info":
		el = l.logger.Info().Str("eventClass", eventClass).Str("event", eventName)
	case "warn":
		el = l.logger.Warn().Str("eventClass", eventClass).Str("event", eventName)
	case "error":
		el = l.logger.Error().Str("eventClass", eventClass).Str("event", eventName)
	case "fatal":
		el = l.logger.Fatal().Str("eventClass", eventClass).Str("event", eventName)
	default:
		el = l.logger.Info().Str("eventClass", eventClass).Str("event", eventName)
	}

	if len(args) == 0 {
		el.Ctx(ctx).Msg(message)
	} else {
		el.Ctx(ctx).Msgf(message, args...)
	}
}

type TracingHook struct{}

func (h TracingHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	ctx := e.GetCtx()
	startTime, ok := ctx.Value(ContextKeyStartTime{}).(time.Time)
	if !ok {
		startTime = time.Now()
	}
	elapsedtime := time.Since(startTime)
	// set base as nanosecond, can be customized to other unit
	nanosecond := float64(elapsedtime.Nanoseconds())
	// convert to millisecond
	millisecond := nanosecond / float64(time.Millisecond)
	e.Float64("elapsedTime", millisecond)
}
