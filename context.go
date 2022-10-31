package capo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tonygcs/capo/marshaler"
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
	err         error
	w           http.ResponseWriter
	r           *http.Request
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
	}
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
	return ctx.ctx.Err()
}

// Value returns any value in the request context.
func (ctx *Context) Value(key any) any {
	return ctx.ctx.Value(key)
}

// Read takes the information in the request body and unmarshal the data in the
// entity provided.
func (ctx *Context) Read(entity interface{}) error {
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
func (ctx *Context) Write(entity interface{}) error {
	data, err := m.Marshal(entity)
	if err != nil {
		return fmt.Errorf("invalid entity format :: %w", err)
	}
	_, err = ctx.w.Write(data)
	if err != nil {
		return fmt.Errorf("cannot write the response :: %w", err)
	}
	return nil
}

// Status sets the response status and returns itself.
func (ctx *Context) Status(status int) *Context {
	ctx.w.WriteHeader(status)
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
