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
	id := c.Request().Header.Get(echo.HeaderXRequestID)
	if id == "" {
		id = c.Response().Header().Get(echo.HeaderXRequestID)
	}
	return id
}

func NewMiddlewareLogger(log *zerolog.Logger) *LoggerMiddleware {
	logger := log.With().
		Str("eventClass", "echo.middleware").
		Str("event", "request").
		Logger()

	return &LoggerMiddleware{
		logger: &logger,
	}
}

func (l *LoggerMiddleware) getContentLength(c echo.Context) int64 {
	cl := c.Request().Header.Get(echo.HeaderContentLength)
	if cl == "" {
		return 0
	}
	clInt, err := strconv.ParseInt(cl, 10, 64)
	if err != nil {
		return 0
	}
	return clInt
}

// LogHandler logs request and response details.
func (l *LoggerMiddleware) LogHandler(c echo.Context, v middleware.RequestLoggerValues) error {
	contentLength := l.getContentLength(c)

	logEvent := l.logger.Info()
	errStr := ""
	if v.Error != nil {
		logEvent = l.logger.Error()
		errStr = v.Error.Error()
	}

	logEvent.
		Str("eventID", v.RequestID).
		Str("remoteIP", c.RealIP()).
		Str("host", c.Request().Host).
		Str("method", c.Request().Method).
		Str("uri", c.Request().RequestURI).
		Str("userAgent", c.Request().UserAgent()).
		Int("status", c.Response().Status).
		Str("requestTime", v.StartTime.Format(time.RFC3339Nano)).
		Int64("latency", v.Latency.Milliseconds()).
		Str("latencyHuman", v.Latency.String()).
		Int64("bytesIn", contentLength).
		Int64("bytesOut", c.Response().Size).
		Msg(errStr)

	return nil
}
