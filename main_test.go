package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServerHandleGetRequest(t *testing.T) {
	serverHandler := New()

	msg := "hello test"
	serverHandler.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(msg))
	})

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	c := &http.Client{}

	r, err := c.Get(s.URL)
	require.NoError(t, err)

	bodyContent, err := io.ReadAll(r.Body)
	require.NoError(t, err)
	require.Equal(t, msg, string(bodyContent))
}

func TestServerNewGroup(t *testing.T) {
	serverHandler := New()

	msg := "hello test"
	g := serverHandler.Group("/")
	g.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(msg))
	})

	s := httptest.NewServer(serverHandler)
	defer s.Close()

	c := &http.Client{}

	r, err := c.Get(s.URL)
	require.NoError(t, err)

	bodyContent, err := io.ReadAll(r.Body)
	require.NoError(t, err)
	require.Equal(t, msg, string(bodyContent))
}
