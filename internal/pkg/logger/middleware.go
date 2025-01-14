package logger

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

// LoggerMiddleware provides logging functionality for Echo requests.
type LoggerMiddleware struct {
	logger *zerolog.Logger
}

// NewLoggerMiddleware creates a new LoggerMiddleware instance with a logger
// initialized to include common fields for all incoming requests.
func NewLoggerMiddleware(log *zerolog.Logger) *LoggerMiddleware {
	logger := log.With().
		Str("eventClass", "echo.middleware").
		Str("event", "request").
		Logger()

	return &LoggerMiddleware{
		logger: &logger,
	}
}

// LogRequest logs detailed information about each incoming request and its response.
// It sets the log level to Info by default, switching to Error if an error occurred.
func (l *LoggerMiddleware) LogRequest(c echo.Context, v middleware.RequestLoggerValues) error {
	logEvent := l.logger.Info()
	errStr := ""

	// If there's an error, switch log level to Error and capture the error message.
	if v.Error != nil {
		logEvent = l.logger.Error()
		errStr = v.Error.Error()
	}

	// Log the request and response details with structured fields.
	logEvent.
		Str("eventID", v.RequestID).
		Str("remoteIP", c.RealIP()).
		Str("host", c.Request().Host).
		Str("method", c.Request().Method).
		Str("uri", c.Request().RequestURI).
		Str("userAgent", c.Request().UserAgent()).
		Int("status", c.Response().Status).
		Str("requestTime", v.StartTime.Format(time.RFC3339Nano)).
		Int64("latency", v.Latency.Nanoseconds()).
		Str("latencyHuman", v.Latency.String()).
		Int64("bytesIn", l.getContentLength(c)).
		Int64("bytesOut", c.Response().Size).
		Msg(errStr)

	return nil
}

// getContentLength retrieves and parses the Content-Length header from the request.
// If parsing fails or the header is missing, it returns 0.
func (l *LoggerMiddleware) getContentLength(c echo.Context) int64 {
	if cl := c.Request().Header.Get(echo.HeaderContentLength); cl != "" {
		if clInt, err := strconv.ParseInt(cl, 10, 64); err == nil {
			return clInt
		}
	}
	return 0
}
