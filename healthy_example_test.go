package capo

import (
	"context"
	"net/http"
)

func ExampleHealthyServer() {
	opts := &Opts{
		// Set the functions here to provide a method that will check the server
		// status.
		Health:    func(ctx context.Context) error { return nil },
		Liveness:  func(ctx context.Context) error { return nil },
		Readiness: func(ctx context.Context) error { return nil },
	}
	hs := NewHealthyServer(opts)

	go http.ListenAndServe(":8008", hs)

	// Run you server here in other port
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}
