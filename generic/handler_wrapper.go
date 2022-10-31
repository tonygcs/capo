package generic

import (
	capo "github.com/tonygcs/capo"
)

// Handler is the function definition for a HTTP handler.
type Handler[T any, U any] func(ctx *Context[T, U]) error

// WrapGenericHandler wraps a HTTP handler function to provide generics.
func WrapGenericHandler[T any, U any](handler Handler[T, U]) capo.Handler {
	return func(ctx *capo.Context) error {
		newCtx := NewContext[T, U](ctx)
		err := newCtx.load()
		if err != nil {
			return err
		}

		return handler(newCtx)
	}
}
