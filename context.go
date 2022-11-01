package capo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tonygcs/capo/marshaler"
	"github.com/tonygcs/gnalog"
)

var (
	// ErrEmptyBody indicates the request body is empty.
	ErrEmptyBody = errors.New("the body is empty")
)

var m = marshaler.GetMarshaler()

// Context is the request context.
type Context struct {
	ctx         context.Context
	cancelled   bool
	cancelCtxFn func()
	w           http.ResponseWriter
	r           *http.Request

	logger gnalog.Logger

	err          error
	status       int
	responseData any
	headers      map[string]string
}

// NewContext creates a new instance of context.
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx, cancel := context.WithCancel(r.Context())
	r = r.WithContext(ctx)

	return &Context{
		ctx:         ctx,
		cancelCtxFn: cancel,
		cancelled:   false,
		r:           r,
		w:           w,
		headers:     make(map[string]string),
	}
}

// Request returns the http request entity.
func (ctx *Context) Request() *http.Request {
	return ctx.r
}

// AddHeader includes a header value in the response.
func (ctx *Context) AddHeader(key string, value string) {
	ctx.headers[key] = value
}

// Deadline is the context deadline.
func (ctx *Context) Deadline() (time.Time, bool) {
	return ctx.ctx.Deadline()
}

// Done is the context done channel.
func (ctx *Context) Done() <-chan struct{} {
	return ctx.ctx.Done()
}

// Err is the context error.
func (ctx *Context) Err() error {
	return ctx.err
}

// Value returns any value in the request context.
func (ctx *Context) Value(key any) any {
	return ctx.ctx.Value(key)
}

// With includes a value in the current context.
func (ctx *Context) With(key any, value any) {
	newCtx := context.WithValue(ctx.ctx, key, value)
	ctx.ctx = newCtx
}

// Read takes the information in the request body and unmarshal the data in the
// entity provided.
func (ctx *Context) Read(entity any) error {
	data, err := io.ReadAll(ctx.r.Body)
	if err != nil {
		return fmt.Errorf("cannot read request body :: %w", err)
	}

	if len(data) == 0 {
		return ErrEmptyBody
	}

	err = m.Unmarshal(data, entity)
	if err != nil {
		return fmt.Errorf("invalid body format :: %w", err)
	}
	return nil
}

// Write marshals and write the information in the entity provided into the http
// response.
func (ctx *Context) Write(entity any) *Context {
	ctx.responseData = entity
	return ctx
}

// Status returns the status code that the server will return to the client.
func (ctx *Context) Status() int {
	if ctx.status <= 0 {
		return http.StatusOK
	}

	return ctx.status
}

// SetStatus sets the response status and returns itself.
func (ctx *Context) SetStatus(status int) *Context {
	ctx.status = status
	return ctx
}

// Cancel sets the error in the context and cancel it.
func (ctx *Context) Cancel(err error) error {
	ctx.err = err

	if !ctx.cancelled {
		ctx.cancelled = true
		ctx.cancelCtxFn()
	}

	return err
}

// Logger returns the logger for the current context.
func (ctx *Context) Logger() gnalog.Logger {
	if ctx.logger != nil {
		// Create default logger if it does not exists.
		ctx.SetLogger(gnalog.New())
	}

	return ctx.logger
}

// SetLogger sets the context logger.
func (ctx *Context) SetLogger(logger gnalog.Logger) {
	ctx.logger = logger
}

func (ctx *Context) closeResponse() error {
	// Set response status code.
	if ctx.status > 0 {
		ctx.w.WriteHeader(ctx.status)
	}

	// Add the headers.
	for key, value := range ctx.headers {
		ctx.w.Header().Add(key, value)
	}

	// Set the body data if it is needed.
	if ctx.responseData != nil {
		data, err := m.Marshal(ctx.responseData)
		if err != nil {
			return fmt.Errorf("invalid entity format :: %w", err)
		}
		_, err = ctx.w.Write(data)
		if err != nil {
			return fmt.Errorf("cannot write the response :: %w", err)
		}
	}

	return nil
}
