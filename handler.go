package capo

import (
	"context"
	"errors"
	"io"
	"net/http"
)

// ContextReader is the function type that the application needs to use to
// extend the context for each request. The 'ctx' argument is the current
// request context that was created inside the wrapper. It already has the
// request id and the logger.
type ContextReader func(ctx context.Context) (context.Context, error)

// HandlerWithCtx creates a 'http.Handler' from a method that will extend the
// context for the request and the function that will handle it. It will also
// read the request and write the response.
//
// The wrapper will handle the errors that are returned by the handler to return
// the error response to the clients. It will return a 500 if the error is not
// an instance that was created with the method 'Error'.
//
// The 'ctx' argument in the handler has the request id and the logger that can
// be used in the request flow.
func HandlerWithCtx[T any, U any](
	ctxFn ContextReader,
	h func(ctx context.Context, req T) (U, error),
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := newContext(r)
		if ctxFn != nil {
			newCtx, err := ctxFn(ctx)
			if err != nil {
				handleError(ctx, err, w)
				return
			}
			ctx = newCtx
		}

		var req T
		// Do not read the request body if the data type for the request is
		// "struct{}".
		if _, ok := any(req).(struct{}); !ok {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				err = Error(err, InvalidRequestBodyCode, InvalidRequestBodyData, http.StatusBadRequest)
				handleError(ctx, err, w)
				return
			}
			err = RequestDecoder(body, &req)
			if err != nil {
				err = Error(err, InvalidRequestBodyCode, InvalidRequestBodyMsg, http.StatusBadRequest)
				handleError(ctx, err, w)
				return
			}
		}

		res, err := h(ctx, req)
		if err != nil {
			handleError(ctx, err, w)
			return
		}

		if _, ok := any(res).(*struct{}); ok {
			// Do not write the response if the data type is "nil".
			w.WriteHeader(http.StatusOK)
			return
		}

		resBody, err := ResponseEncoder(res)
		if err != nil {
			handleError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(resBody)
		if err != nil {
			Logger(ctx).Info("cannot write the response in the http connection :: %v", err)
			// Do not return anything because the connection is not available
			// any more.
			return
		}
	})
}

// Handler creates a 'http.Handler' from a function that will handle the
// request. It will read the request and write the response.
//
// The wrapper will handle the errors that are returned by the handler to return
// the error response to the clients. It will return a 500 if the error is not
// an instance that was created with the method 'Error'.
//
// The 'ctx' argument in the handler has the request id and the logger that can
// be used in the request flow.
func Handler[T any, U any](
	h func(ctx context.Context, req T) (U, error),
) http.Handler {
	return HandlerWithCtx(nil, h)
}

func handleError(ctx context.Context, err error, w http.ResponseWriter) {
	var capoErr *capoError
	if !errors.As(err, &capoErr) {
		// Any other error that is not *capoError will be consider as an
		// internal error.
		capoErr = Error(err, InternalErrorCode, InternalErrorMsg, http.StatusInternalServerError)
	}

	resBody, err := ResponseEncoder(capoErr)
	if err != nil {
		// Log an error because the application cannot build its own response.
		Logger(ctx).Error("cannot build the error response :: %v", err)
	}

	w.WriteHeader(capoErr.status)

	if len(resBody) > 0 {
		_, err := w.Write(resBody)
		if err != nil {
			Logger(ctx).Info("cannot write the error response in the http connection :: %v", err)
			// Do not return anything because the connection is not available
			// any more.
		}
	}
}
