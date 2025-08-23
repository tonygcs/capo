package capo

import (
	"context"
	"net/http"
)

func ExampleHandler() {
	get := func(ctx context.Context, _ struct{}) ([]*Response, error) {
		// Use "struct{}" in the arguments that the method will get to indicate
		// that the body is empty.

		// Read values from the context.
		logger := Logger(ctx)
		reqId := ReqId(ctx)

		// Note that the logger already contains the request id in a field.
		logger.Debug("Request ID: %s", reqId)

		return []*Response{
			{Surname: "Surname"},
		}, nil
	}
	post := func(ctx context.Context, req *Request) (*Response, error) {
		return &Response{Surname: "Surname"}, nil
	}
	delete := func(ctx context.Context, req *Request) (*struct{}, error) {
		// Use "*struct{}" in the arguments that the method will return to
		// indicate that the response body is empty.
		return nil, nil
	}

	s := http.NewServeMux()
	s.Handle("GET /", Handler(get))
	s.Handle("POST /", Handler(post))
	s.Handle("DELETE /", Handler(delete))
}

func ExampleHandlerWithCtx() {
	type ctxKey string
	const surnameKey ctxKey = "surnameKey"

	extendCtx := func(ctx context.Context) (context.Context, error) {
		surname := "Surname"
		newCtx := context.WithValue(ctx, surnameKey, surname)
		return newCtx, nil
	}

	h := func(ctx context.Context, req *Request) (*Response, error) {
		surnameCtx := ctx.Value(surnameKey)
		surname, _ := surnameCtx.(string) // Handle the casting here
		return &Response{Surname: surname}, nil
	}

	s := http.NewServeMux()
	s.Handle("GET /", HandlerWithCtx(extendCtx, h))
}
