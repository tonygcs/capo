package generic

import (
	"errors"
	"net/http"
	"time"

	capo "github.com/tonygcs/capo"
	"github.com/tonygcs/gnalog"
)

// Context is the request context.
type Context[T any, U any] struct {
	ctx  *capo.Context
	Data *T
}

// NewContext creates a new instance of a context.
func NewContext[T any, U any](ctx *capo.Context) *Context[T, U] {
	return &Context[T, U]{
		ctx: ctx,
	}
}

// Request returns the http request entity.
func (ctx *Context[T, U]) Request() *http.Request {
	return ctx.ctx.Request()
}

// AddHeader includes a header value in the response.
func (ctx *Context[T, U]) AddHeader(key string, value string) {
	ctx.ctx.AddHeader(key, value)
}

// Deadline is the context deadline.
func (ctx *Context[T, U]) Deadline() (time.Time, bool) {
	return ctx.ctx.Deadline()
}

// Done is the context done channel.
func (ctx *Context[T, U]) Done() <-chan struct{} {
	return ctx.ctx.Done()
}

// Err is the context error.
func (ctx *Context[T, U]) Err() error {
	return ctx.ctx.Err()
}

// Value returns any value in the request context.
func (ctx *Context[T, U]) Value(key any) any {
	return ctx.ctx.Value(key)
}

// With includes a value in the current context.
func (ctx *Context[T, U]) With(key any, value any) {
	ctx.ctx.With(key, value)
}

// Write marshals and write the information in the entity provided into the http
// response.
func (ctx *Context[T, U]) Write(entity *U) *Context[T, U] {
	ctx.ctx.Write(entity)
	return ctx
}

// Cancel sets the error in the context and cancel it.
func (ctx *Context[T, U]) Cancel(err error) {
	ctx.ctx.Cancel(err)
}

// Logger returns the logger for the current context.
func (ctx *Context[T, U]) Logger() gnalog.Logger {
	return ctx.ctx.Logger()
}

// SetLogger sets the context logger.
func (ctx *Context[T, U]) SetLogger(logger gnalog.Logger) {
	ctx.ctx.SetLogger(logger)
}

// load takes the information in the request body and sets the 'Data' field in
// the current context.
func (ctx *Context[T, U]) load() error {
	entity := new(T)

	err := ctx.ctx.Read(entity)
	if errors.Is(err, capo.ErrEmptyBody) {
		entity = nil
	} else if err != nil {
		return err
	}

	ctx.Data = entity
	return nil
}
