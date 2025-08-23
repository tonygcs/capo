# CAPO

Golang HTTP server with health checks and a generic wrapper for the handlers.

## Health checks

The `HealthyServer` provides the three endpoints for the health checks:
liveness, readiness and health. See the file
[healthy_example_test.go](.healthy_example_test.go) to check how the
`HealthyServer` is used.

## Generic handler

It's a wrapper to add strong typing and error handling to every http handler. It
provides a context with a request ID and logger to be used in the handler
method. See the file [handler_example_test.go](.handler_example_test.go) to
review the `Handler` and `HandlerWithCtx` usage.
