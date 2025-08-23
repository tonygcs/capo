package capo

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/tonygcs/gnalog"
)

type ctxKey string

const reqIdKey ctxKey = "reqId"
const loggerKey ctxKey = "logger"
const reqKey ctxKey = "request"

// newContext extends the request context and include the request id and the
// logger on it.
//
// It also includes the '*http.Request'.
func newContext(r *http.Request) context.Context {
	ctx := r.Context()

	reqId := uuid.NewString()
	ctx = context.WithValue(ctx, reqIdKey, reqId)

	logger := LoggerFactory()
	logger = logger.With("reqId", reqId)
	ctx = context.WithValue(ctx, loggerKey, logger)

	r = r.WithContext(ctx)
	ctx = context.WithValue(ctx, reqKey, r)

	return ctx
}

// ReqId returns the request id for the current request.
func ReqId(ctx context.Context) string {
	value := ctx.Value(reqIdKey)
	if reqId, ok := value.(string); ok {
		return reqId
	}
	panic("the context does not have a request id")
}

// Req returns the '*http.Request' instance for the current request.
func Req(ctx context.Context) *http.Request {
	value := ctx.Value(reqKey)
	if req, ok := value.(*http.Request); ok {
		return req
	}
	panic("the context does not have a request")
}

// Logger returns the logger with all the context for the current request.
func Logger(ctx context.Context) gnalog.Logger {
	value := ctx.Value(loggerKey)
	if logger, ok := value.(gnalog.Logger); ok {
		return logger
	}
	panic("the context does not have a logger")
}
