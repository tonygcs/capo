package generic

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tonygcs/capo"
)

type TestEntity struct {
	Message string `json:"msg"`
}

func TestWrapGenericHandlerCanSetResponseBody(t *testing.T) {
	h := capo.New()

	msg := "test generic message"
	h.Get("", WrapGenericHandler(func(ctx *Context[any, TestEntity]) error {
		ctx.Write(&TestEntity{Message: msg})
		return nil
	}))

	s := httptest.NewServer(h)
	defer s.Close()

	res, err := NewRequest[any, TestEntity]().URL(s.URL).Do()
	require.NoError(t, err)
	require.Equal(t, msg, res.Message)
}

func TestWrapGenericHandlerCanReadRequestBody(t *testing.T) {
	h := capo.New()

	msg := "test generic message"
	h.Post("", WrapGenericHandler(func(ctx *Context[TestEntity, any]) error {
		require.Equal(t, msg, ctx.Data.Message)
		return nil
	}))

	s := httptest.NewServer(h)
	defer s.Close()

	_, err := NewRequest[TestEntity, any]().URL(s.URL).Method(http.MethodPost).Data(&TestEntity{Message: msg}).Do()
	require.NoError(t, err)
}
