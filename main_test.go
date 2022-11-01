package capo

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tonygcs/capo/client"
)

type TestData struct {
	Message string `json:"msg"`
}

func TestServerImplementsGroup(t *testing.T) {
	var _ Group = New()
}

func TestServerHandleGetRequest(t *testing.T) {
	serverHandler := New()

	msg := "hello test"
	serverHandler.Get("/", func(ctx *Context) error {
		ctx.Write(&TestData{Message: msg})
		return nil
	})

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	req := client.NewRequest().URL(s.URL)
	res := &TestData{}
	err := req.Do(res)
	require.NoError(t, err)
	require.Equal(t, msg, res.Message)
}

func TestServerHandlePostRequest(t *testing.T) {
	serverHandler := New()

	reqMsg := "hello request test"
	resMsg := "hello response test"
	serverHandler.Post("/post", func(ctx *Context) error {
		// Check request data.
		req := &TestData{}
		err := ctx.Read(req)
		require.NoError(t, err)
		require.Equal(t, reqMsg, req.Message)

		ctx.Write(&TestData{Message: resMsg})
		return nil
	})

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	req := client.NewRequest().URL(s.URL).RelativePath("post").Method(http.MethodPost).Data(&TestData{Message: reqMsg})
	res := &TestData{}
	err := req.Do(res)
	require.NoError(t, err)
	require.Equal(t, resMsg, res.Message)
}

func TestServerHandlePutRequest(t *testing.T) {
	serverHandler := New()

	reqMsg := "hello request test"
	resMsg := "hello response test"
	serverHandler.Put("/put", func(ctx *Context) error {
		// Check request data.
		req := &TestData{}
		err := ctx.Read(req)
		require.NoError(t, err)
		require.Equal(t, reqMsg, req.Message)

		ctx.Write(&TestData{Message: resMsg})
		return nil
	})

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	req := client.NewRequest().URL(s.URL).RelativePath("put").Method(http.MethodPut).Data(&TestData{Message: reqMsg})
	res := &TestData{}
	err := req.Do(res)
	require.NoError(t, err)
	require.Equal(t, resMsg, res.Message)
}

func TestServerHandleDeleteRequest(t *testing.T) {
	serverHandler := New()

	msg := "hello delete test"
	serverHandler.Delete("/delete", func(ctx *Context) error {
		ctx.Write(&TestData{Message: msg})
		return nil
	})

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	req := client.NewRequest().URL(s.URL).RelativePath("delete").Method(http.MethodDelete)
	res := &TestData{}
	err := req.Do(res)
	require.NoError(t, err)
	require.Equal(t, msg, res.Message)
}

func TestRequestRunsAllMiddlewareTypes(t *testing.T) {
	serverHandler := New()
	calls := 0
	serverHandler.UseBefore(func(ctx *Context) error { calls++; return nil })
	serverHandler.UseAfter(func(ctx *Context) error { calls++; return nil })
	serverHandler.UseAfterAlways(func(ctx *Context) { calls++ })

	serverHandler.Get("", func(ctx *Context) error { calls++; return nil })

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	err := client.NewRequest().URL(s.URL).Do(nil)
	require.NoError(t, err)
	require.Equal(t, 4, calls)
}

func TestHandlerErrorStopsAfterMiddlewares(t *testing.T) {
	serverHandler := New()
	calls := 0
	serverHandler.UseBefore(func(ctx *Context) error { calls++; return nil })
	serverHandler.UseAfter(func(ctx *Context) error { calls++; return nil })

	serverHandler.Get("", func(ctx *Context) error { calls++; return errors.New("test error") })

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	err := client.NewRequest().URL(s.URL).Do(nil)
	require.NoError(t, err)
	require.Equal(t, 2, calls)
}

func TestHandlerErrorRunsAfterAlwaysMiddlewares(t *testing.T) {
	serverHandler := New()
	calls := 0
	serverHandler.UseBefore(func(ctx *Context) error { calls++; return nil })
	serverHandler.UseAfter(func(ctx *Context) error { calls++; return nil })
	serverHandler.UseAfterAlways(func(ctx *Context) { calls++ })

	serverHandler.Get("", func(ctx *Context) error { calls++; return errors.New("test error") })

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	err := client.NewRequest().URL(s.URL).Do(nil)
	require.NoError(t, err)
	require.Equal(t, 3, calls)
}

func TestBeforeMiddlewareStopsPropagation(t *testing.T) {
	serverHandler := New()
	calls := 0
	serverHandler.UseBefore(func(ctx *Context) error { calls++; return errors.New("test error") })
	serverHandler.UseAfter(func(ctx *Context) error { calls++; return nil })

	serverHandler.Get("", func(ctx *Context) error { calls++; return nil })

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	err := client.NewRequest().URL(s.URL).Do(nil)
	require.NoError(t, err)
	require.Equal(t, 1, calls)
}

func TestAfterMiddlewareStopsPropagation(t *testing.T) {
	serverHandler := New()
	calls := 0
	serverHandler.UseBefore(func(ctx *Context) error { calls++; return nil })
	serverHandler.UseAfter(func(ctx *Context) error { calls++; return errors.New("test error") })
	serverHandler.UseAfter(func(ctx *Context) error { calls++; return nil })

	serverHandler.Get("", func(ctx *Context) error { calls++; return nil })

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	err := client.NewRequest().URL(s.URL).Do(nil)
	require.NoError(t, err)
	require.Equal(t, 3, calls)
}

func TestAfterAlwaysIsCalledIfRequestPanics(t *testing.T) {
	h := New()
	calls := 0
	h.UseAfterAlways(func(ctx *Context) { calls++ })

	h.Get("", func(ctx *Context) error { panic("error") })

	s := httptest.NewServer(h)
	defer s.Close()

	err := client.NewRequest().URL(s.URL).Do(nil)
	require.NoError(t, err)
	require.Equal(t, 1, calls)
}
