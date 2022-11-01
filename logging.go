package capo

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tonygcs/gnalog"
)

type contextKey string

const (
	defaultReqIDHeaderKey string     = "X-Request-ID"
	incomingTimeKey       contextKey = "INCOMING_TIME_KEY"
)

// CreateLog returns the middleware to create a log entity in the request
// context. It also includes the request id in the logger context and the
// response header.
func CreateLog(requestIDHeaderKey string) Handler {
	if requestIDHeaderKey == "" {
		requestIDHeaderKey = defaultReqIDHeaderKey
	}

	return func(ctx *Context) error {
		var l gnalog.Logger = gnalog.New()

		requestID := uuid.New().String()
		l = l.With("request-id", requestID)

		ctx.AddHeader(requestIDHeaderKey, requestID)

		ctx.SetLogger(l)
		ctx.With(incomingTimeKey, time.Now())

		return nil
	}
}

// LogRequest logs the request that is about to finish. It logs (INFO level) the
// request result status, method, endpoint, error (if it exists) and run time.
// In case the result status is an internal error, the middleware will log the
// record as an ERROR.
func LogRequest(ctx *Context) {
	status := ctx.Status()

	l := ctx.Logger().
		With("method", ctx.r.Method).
		With("endpoint", ctx.r.URL.Path).
		With("status", status)

	// Log the error message.
	if err := ctx.Err(); err != nil {
		l.With("error", err.Error())
	}

	// Calculate time.
	start := ctx.Value(incomingTimeKey).(time.Time)
	delay := time.Since(start).Milliseconds()
	l = l.With("run-time", delay)

	// Log an error if the response status is not valid.
	if status >= http.StatusInternalServerError {
		l.Error("internal error")
	} else {
		l.Info("handled error")
	}
}
