package logger

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

type LoggerMiddleware struct {
	logger *zerolog.Logger
}

func getRequestID(c echo.Context) string {
	// Fetch X-Request-ID from request header, fallback to response header if not present
	if reqID := c.Request().Header.Get(echo.HeaderXRequestID); reqID != "" {
		return reqID
	}
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

func NewLoggerMiddleware(log *zerolog.Logger) *LoggerMiddleware {
	// Initialize logger with common fields for all requests
	logger := log.With().
		Str("eventClass", "echo.middleware").
		Str("event", "request").
		Logger()

	return &LoggerMiddleware{
		logger: &logger,
	}
}

func (l *LoggerMiddleware) getContentLength(c echo.Context) int64 {
	// Parse content length from request header
	if cl := c.Request().Header.Get(echo.HeaderContentLength); cl != "" {
		if clInt, err := strconv.ParseInt(cl, 10, 64); err == nil {
			return clInt
		}
	}
	return 0
}

// LogRequest logs detailed information about each request and response
func (l *LoggerMiddleware) LogRequest(c echo.Context, v middleware.RequestLoggerValues) error {
	// Set default log level to info, switch to error if there's an error
	logEvent := l.logger.Info()
	errStr := ""
	if v.Error != nil {
		logEvent = l.logger.Error()
		errStr = v.Error.Error()
	}

	// Log the request with detailed, descriptive fields
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
