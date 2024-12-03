package logger

import (
	"bytes"
	"encoding/json"
	"net/http/httputil"
	"strconv"
	"time"

	"golang-boilerplate/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

const XHEADERREQUESTIME = "X-Custom-RequestTime"

// Logger -.
type MiddlewareLogger struct {
	logger *zerolog.Logger
}

// ServiceRequestTime middleware adds a `X-Custom-RequestTime` header to the Request for logging.
func ServiceRequestTime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().Header.Set("X-Custom-RequestTime", time.Now().Format(time.RFC3339Nano))
		return next(c)
	}
}

func getEchoRequestID(ctx echo.Context) string {
	req := ctx.Request()
	res := ctx.Response()
	id := req.Header.Get(echo.HeaderXRequestID)
	if id == "" {
		id = res.Header().Get(echo.HeaderXRequestID)
	}
	return id
}

// New -.
func NewMiddlewareLogger(_ *config.Logger, log *zerolog.Logger) *MiddlewareLogger {
	logger := log.With().
		Str("eventClass", "echo.middleware").
		Str("event", "request").
		Logger()

	return &MiddlewareLogger{
		logger: &logger,
	}
}

func (l *MiddlewareLogger) getContentLength(c echo.Context) int64 {
	// Content Length
	cl := c.Request().Header.Get(echo.HeaderContentLength)
	if cl == "" {
		cl = "0"
	}
	clint, err := strconv.ParseInt(cl, 10, 64)
	if err != nil {
		clint = 0
	}
	return clint
}

func (l *MiddlewareLogger) getElapsedTimeFromCustomHeader(c echo.Context) time.Duration {
	reqTime, errT := time.Parse(time.RFC3339, c.Request().Header.Get(XHEADERREQUESTIME))
	var elapstime time.Duration
	if errT == nil {
		elapstime = time.Since(reqTime)
	}
	return elapstime
}

// MiddlewareLogHandler -.
func (l *MiddlewareLogger) MiddlewareLogHandler(c echo.Context, v middleware.RequestLoggerValues) error {
	clint := l.getContentLength(c)

	errStr := ""
	el := l.logger.Info()
	if v.Error != nil {
		el = l.logger.Error()
		b, err := json.Marshal(v.Error.Error())
		if err != nil {
			errStr = v.Error.Error()
		} else {
			b = b[1 : len(b)-1]
			errStr = string(b)
		}
	}

	el.
		Str("eventID", v.RequestID).
		Str("remoteIP", c.RealIP()).
		Str("host", c.Request().Host).
		Str("method", c.Request().Method).
		Str("uri", c.Request().RequestURI).
		Str("userAgent", c.Request().UserAgent()).
		Int("status", c.Response().Status).
		Str("requestTime", c.Request().Header.Get(XHEADERREQUESTIME)).
		Int64("latency", int64(v.Latency)).
		Str("latencyHuman", v.Latency.String()).
		Int64("bytesIn", clint).
		Int64("bytesOut", c.Response().Size).
		Msg(errStr)

	return nil
}

// MiddlewareBodyDumpLogHandler -.
// Only Log Info level with request response body dump
// No Error Log like default logger.
func (l *MiddlewareLogger) MiddlewareBodyDumpLogHandler(c echo.Context, _, res []byte) {
	// Get elapsed time
	elapstime := l.getElapsedTimeFromCustomHeader(c)

	// Content Length
	clint := l.getContentLength(c)

	var b bytes.Buffer
	var respHeader []byte
	err := c.Response().Header().Write(&b)
	if err != nil {
		respHeader, err = json.Marshal(c.Response().Header())
		if err != nil {
			respHeader = []byte("Error Marshal Response Header : " + err.Error())
		}
	} else {
		respHeader = b.Bytes()
	}

	msg := ""
	requestDump, err := httputil.DumpRequest(c.Request(), true)
	if err != nil {
		msg += "Error Dump Request : " + err.Error()
	}

	l.logger.Info().
		Str("eventID", getEchoRequestID(c)).
		Str("remoteIP", c.RealIP()).
		Str("host", c.Request().Host).
		Str("method", c.Request().Method).
		Str("uri", c.Request().RequestURI).
		Str("userAgent", c.Request().UserAgent()).
		Int("status", c.Response().Status).
		Str("requestTime", c.Request().Header.Get(XHEADERREQUESTIME)).
		Int64("latency", int64(elapstime)).
		Str("latencyHuman", elapstime.String()).
		Int64("bytesIn", clint).
		Int64("bytesOut", c.Response().Size).
		Str("request", string(requestDump)).
		Str("response", string(respHeader)+"\r\n"+string(res)).
		Msg(msg)
}
