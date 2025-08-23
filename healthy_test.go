package capo

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func doRequest(handler http.Handler, path string) int {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec.Code
}

func TestLiveness_NoOptDefined_Returns200(t *testing.T) {
	server := NewHealthyServer(nil)

	code := doRequest(server, "/liveness")
	if code != http.StatusOK {
		t.Errorf("expected 200, got %d", code)
	}

	code = doRequest(server, "/readiness")
	if code != http.StatusOK {
		t.Errorf("expected 200, got %d", code)
	}

	code = doRequest(server, "/health")
	if code != http.StatusOK {
		t.Errorf("expected 200, got %d", code)
	}
}

func TestLiveness_NoFuncDefined_Returns200(t *testing.T) {
	server := NewHealthyServer(&Opts{})
	code := doRequest(server, "/liveness")

	if code != http.StatusOK {
		t.Errorf("expected 200, got %d", code)
	}
}

func TestLiveness_FuncDefined_OK_Returns200(t *testing.T) {
	server := NewHealthyServer(&Opts{
		Liveness: func(ctx context.Context) error { return nil },
	})
	code := doRequest(server, "/liveness")

	if code != http.StatusOK {
		t.Errorf("expected 200, got %d", code)
	}
}

func TestLiveness_FuncDefined_Error_Returns503(t *testing.T) {
	server := NewHealthyServer(&Opts{
		Liveness: func(ctx context.Context) error { return errors.New("fail") },
	})
	code := doRequest(server, "/liveness")

	if code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", code)
	}
}

func TestReadiness_NoFuncDefined_Returns200(t *testing.T) {
	server := NewHealthyServer(&Opts{})
	code := doRequest(server, "/readiness")

	if code != http.StatusOK {
		t.Errorf("expected 200, got %d", code)
	}
}

func TestReadiness_FuncDefined_OK_Returns200(t *testing.T) {
	server := NewHealthyServer(&Opts{
		Readiness: func(ctx context.Context) error { return nil },
	})
	code := doRequest(server, "/readiness")

	if code != http.StatusOK {
		t.Errorf("expected 200, got %d", code)
	}
}

func TestReadiness_FuncDefined_Error_Returns503(t *testing.T) {
	server := NewHealthyServer(&Opts{
		Readiness: func(ctx context.Context) error { return errors.New("not ready") },
	})
	code := doRequest(server, "/readiness")

	if code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", code)
	}
}

func TestHealthCheck_NoFuncDefined_Returns200(t *testing.T) {
	server := NewHealthyServer(&Opts{})
	code := doRequest(server, "/health")

	if code != http.StatusOK {
		t.Errorf("expected 200, got %d", code)
	}
}

func TestHealthCheck_FuncDefined_OK_Returns200(t *testing.T) {
	server := NewHealthyServer(&Opts{
		Health: func(ctx context.Context) error { return nil },
	})
	code := doRequest(server, "/health")

	if code != http.StatusOK {
		t.Errorf("expected 200, got %d", code)
	}
}

func TestHealthCheck_FuncDefined_Error_Returns503(t *testing.T) {
	server := NewHealthyServer(&Opts{
		Health: func(ctx context.Context) error { return errors.New("unhealthy") },
	})
	code := doRequest(server, "/health")

	if code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", code)
	}
}
