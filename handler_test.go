package capo

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Request struct {
	Name string `json:"name"`
}

type Response struct {
	Surname string `json:"surname"`
}

var validCtxFn = func(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

var validHandler = func(ctx context.Context, req *Request) (*Response, error) {
	return &Response{Surname: "ValidSurname"}, nil
}

type errReader struct{}

func (errReader) Read(p []byte) (int, error) { return 0, errors.New("boom") }
func (errReader) Close() error               { return nil }

func TestHandle_ContextFnError(t *testing.T) {
	ctxFn := func(ctx context.Context) (context.Context, error) {
		err := errors.New("ctx error")
		return nil, Error(err, "BAD_REQ", "Bad request", http.StatusBadRequest)
	}
	handler := HandlerWithCtx(ctxFn, validHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Should return http.StatusBadRequest because ctxFn fails with that error.
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 500, got %d", rr.Code)
	}
}

func TestHandle_InvalidBody(t *testing.T) {
	handler := HandlerWithCtx(validCtxFn, validHandler)

	// Body with a reader that errors on Read.
	badBody := io.NopCloser(errReader{})

	req := httptest.NewRequest(http.MethodPost, "/", badBody)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestHandle_InvalidJSONRequest(t *testing.T) {
	// Patch request decoder.
	oldRequestDecoder := RequestDecoder
	RequestDecoder = func(data []byte, v any) error { return errors.New("bad json") }
	defer func() {
		RequestDecoder = oldRequestDecoder
	}()

	handler := HandlerWithCtx(validCtxFn, validHandler)

	req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewBufferString("{}")))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestHandle_HandlerError(t *testing.T) {
	h := func(ctx context.Context, req map[string]string) (string, error) {
		err := errors.New("handler error")
		return "", Error(err, "BAD_REQ", "Bad request", http.StatusBadRequest)
	}
	handler := HandlerWithCtx(validCtxFn, h)

	req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewBufferString("{}")))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestHandle_InvalidJSONResponse(t *testing.T) {
	// Patch response encoder.
	oldResponseEncoder := ResponseEncoder
	ResponseEncoder = func(data any) ([]byte, error) {
		err := errors.New("handler error")
		return nil, Error(err, "BAD_REQ", "Bad request", http.StatusBadRequest)
	}
	defer func() {
		ResponseEncoder = oldResponseEncoder
	}()

	handler := HandlerWithCtx(validCtxFn, validHandler)

	req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewBufferString("{}")))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestHandle_InternalErrorIfUnknownError(t *testing.T) {
	h := func(ctx context.Context, req *Request) (*Response, error) {
		return nil, errors.New("invalid error")
	}
	handler := HandlerWithCtx(validCtxFn, h)

	req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewBufferString("{}")))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}
}

func TestHandle_RequestStructSkipBody(t *testing.T) {
	// request type is struct{} so body is ignored
	h := func(ctx context.Context, req struct{}) (string, error) {
		return "ok", nil
	}
	handler := HandlerWithCtx(validCtxFn, h)

	// Body with a reader that errors on Read. If the test passes, it means the
	// handle didn't read the request.
	badBody := io.NopCloser(errReader{})

	req := httptest.NewRequest(http.MethodPost, "/", badBody)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestHandle_ResponsePointerStructSkipWrite(t *testing.T) {
	// Patch response encoder so it returns an error. If the response is a 200,
	// it means the request didn't write the response.
	oldResponseEncoder := ResponseEncoder
	ResponseEncoder = func(data any) ([]byte, error) {
		err := errors.New("handler error")
		return nil, Error(err, "BAD_REQ", "Bad request", http.StatusBadRequest)
	}
	defer func() {
		ResponseEncoder = oldResponseEncoder
	}()

	h := func(ctx context.Context, req *Request) (*struct{}, error) {
		return nil, nil
	}
	handler := HandlerWithCtx(validCtxFn, h)

	req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewBufferString("{}")))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestHandle_Success(t *testing.T) {
	handler := HandlerWithCtx(validCtxFn, validHandler)

	req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewBufferString("{}")))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if rr.Body.String() != "{\"surname\":\"ValidSurname\"}\n" {
		t.Errorf("unexpected body, got %s", rr.Body.String())
	}
}
