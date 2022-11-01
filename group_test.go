package capo

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tonygcs/capo/client"
)

func TestGroupHandleGetRequest(t *testing.T) {
	serverHandler := New()
	group := serverHandler.Group("")

	msg := "get test message"
	group.Get("/", func(ctx *Context) error {
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

func TestGroupHandlePostRequest(t *testing.T) {
	serverHandler := New()
	group := serverHandler.Group("")

	reqMsg := "post test request"
	resMsg := "post test response"
	group.Post("/post", func(ctx *Context) error {
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

func TestGroupHandlePutRequest(t *testing.T) {
	serverHandler := New()
	group := serverHandler.Group("put")

	reqMsg := "post test request"
	resMsg := "post test response"
	group.Put("/", func(ctx *Context) error {
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

func TestGroupHandleDeleteRequest(t *testing.T) {
	serverHandler := New()
	group := serverHandler.Group("delete")

	msg := "delete request test"
	group.Delete("/request", func(ctx *Context) error {
		ctx.Write(&TestData{Message: msg})
		return nil
	})

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	req := client.NewRequest().URL(s.URL).RelativePath("delete/request").Method(http.MethodDelete)
	res := &TestData{}
	err := req.Do(res)
	require.NoError(t, err)
	require.Equal(t, msg, res.Message)
}

func TestRequestGroupAndEndpointWithoutPath(t *testing.T) {
	serverHandler := New()
	group := serverHandler.Group("")

	msg := "test request"
	group.Get("", func(ctx *Context) error {
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

func TestRequestGroupWithoutPathAndEndpointWithPath(t *testing.T) {
	serverHandler := New()
	group := serverHandler.Group("")

	msg := "test request"
	group.Get("endpoint-path", func(ctx *Context) error {
		ctx.Write(&TestData{Message: msg})
		return nil
	})

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	req := client.NewRequest().URL(s.URL).RelativePath("endpoint-path")
	res := &TestData{}
	err := req.Do(res)
	require.NoError(t, err)
	require.Equal(t, msg, res.Message)
}

func TestRequestGroupWithPathAndEndpointWithoutPath(t *testing.T) {
	serverHandler := New()
	group := serverHandler.Group("group-path")

	msg := "test request"
	group.Get("", func(ctx *Context) error {
		ctx.Write(&TestData{Message: msg})
		return nil
	})

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	req := client.NewRequest().URL(s.URL).RelativePath("group-path")
	res := &TestData{}
	err := req.Do(res)
	require.NoError(t, err)
	require.Equal(t, msg, res.Message)
}

func TestRequestGroupAndEndpointWithPath(t *testing.T) {
	serverHandler := New()
	group := serverHandler.Group("group-path")

	msg := "test request"
	group.Get("endpoint-path", func(ctx *Context) error {
		ctx.Write(&TestData{Message: msg})
		return nil
	})

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	req := client.NewRequest().URL(s.URL).RelativePath("group-path/endpoint-path")
	res := &TestData{}
	err := req.Do(res)
	require.NoError(t, err)
	require.Equal(t, msg, res.Message)
}

func TestRequestWith5LevelsOfGroups(t *testing.T)  { testRequestWithNLevelsOfGroups(t, 5) }
func TestRequestWith10LevelsOfGroups(t *testing.T) { testRequestWithNLevelsOfGroups(t, 10) }
func TestRequestWith15LevelsOfGroups(t *testing.T) { testRequestWithNLevelsOfGroups(t, 15) }
func TestRequestWith20LevelsOfGroups(t *testing.T) { testRequestWithNLevelsOfGroups(t, 20) }

func testRequestWithNLevelsOfGroups(t *testing.T, levels int) {
	serverHandler := New()
	reqPath := "group-path"
	group := serverHandler.Group(reqPath)

	for i := 0; i < levels; i++ {
		groupPath := fmt.Sprintf("group-path-%d", i)
		group = group.Group(groupPath)
		reqPath = joinPaths(reqPath, groupPath)
	}

	msg := "test request"
	group.Get("", func(ctx *Context) error {
		ctx.Write(&TestData{Message: msg})
		return nil
	})

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	req := client.NewRequest().URL(s.URL).RelativePath(reqPath)
	res := &TestData{}
	err := req.Do(res)
	require.NoError(t, err)
	require.Equal(t, msg, res.Message)
}

func TestGroupRequestRunsAllMiddlewareTypes(t *testing.T) {
	serverHandler := New()
	group := serverHandler.Group("")
	calls := 0
	group.UseBefore(func(ctx *Context) error { calls++; return nil })
	group.UseAfter(func(ctx *Context) error { calls++; return nil })
	group.UseAfterAlways(func(ctx *Context) { calls++ })

	group.Get("", func(ctx *Context) error { calls++; return nil })

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	err := client.NewRequest().URL(s.URL).Do(nil)
	require.NoError(t, err)
	require.Equal(t, 4, calls)
}

func TestGroupHandlerErrorStopsAfterMiddlewares(t *testing.T) {
	serverHandler := New()
	group := serverHandler.Group("")
	calls := 0
	group.UseBefore(func(ctx *Context) error { calls++; return nil })
	group.UseAfter(func(ctx *Context) error { calls++; return nil })

	group.Get("", func(ctx *Context) error { calls++; return errors.New("test error") })

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	err := client.NewRequest().URL(s.URL).Do(nil)
	require.NoError(t, err)
	require.Equal(t, 2, calls)
}

func TestGroupHandlerErrorRunsAfterAlwaysMiddlewares(t *testing.T) {
	serverHandler := New()
	group := serverHandler.Group("")
	calls := 0
	group.UseBefore(func(ctx *Context) error { calls++; return nil })
	group.UseAfter(func(ctx *Context) error { calls++; return nil })
	group.UseAfterAlways(func(ctx *Context) { calls++ })

	group.Get("", func(ctx *Context) error { calls++; return errors.New("test error") })

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	err := client.NewRequest().URL(s.URL).Do(nil)
	require.NoError(t, err)
	require.Equal(t, 3, calls)
}

func TestGroupBeforeMiddlewareStopsPropagation(t *testing.T) {
	serverHandler := New()
	group := serverHandler.Group("")
	calls := 0
	group.UseBefore(func(ctx *Context) error { calls++; return errors.New("test error") })
	group.UseAfter(func(ctx *Context) error { calls++; return nil })

	group.Get("", func(ctx *Context) error { calls++; return nil })

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	err := client.NewRequest().URL(s.URL).Do(nil)
	require.NoError(t, err)
	require.Equal(t, 1, calls)
}

func TestGroupAfterMiddlewareStopsPropagation(t *testing.T) {
	serverHandler := New()
	group := serverHandler.Group("")
	calls := 0
	group.UseBefore(func(ctx *Context) error { calls++; return nil })
	group.UseAfter(func(ctx *Context) error { calls++; return errors.New("test error") })
	group.UseAfter(func(ctx *Context) error { calls++; return nil })

	group.Get("", func(ctx *Context) error { calls++; return nil })

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	err := client.NewRequest().URL(s.URL).Do(nil)
	require.NoError(t, err)
	require.Equal(t, 3, calls)
}

func TestBeforeErrorOnRootStopsGroupBeforeMiddleware(t *testing.T) {
	serverHandler := New()
	calls := 0
	serverHandler.UseBefore(func(ctx *Context) error { calls++; return errors.New("test error") })

	group := serverHandler.Group("")
	group.UseBefore(func(ctx *Context) error { calls++; return nil })
	group.UseAfter(func(ctx *Context) error { calls++; return nil })

	group.Get("", func(ctx *Context) error { calls++; return nil })

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	err := client.NewRequest().URL(s.URL).Do(nil)
	require.NoError(t, err)
	require.Equal(t, 1, calls)
}

func TestAfterErrorOnRootSopsGroupAfterMiddleware(t *testing.T) {
	serverHandler := New()
	calls := 0
	serverHandler.UseAfter(func(ctx *Context) error { calls++; return errors.New("test error") })

	group := serverHandler.Group("")
	group.UseAfter(func(ctx *Context) error { calls++; return nil })

	group.Get("", func(ctx *Context) error { calls++; return nil })

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	err := client.NewRequest().URL(s.URL).Do(nil)
	require.NoError(t, err)
	require.Equal(t, 2, calls)
}
