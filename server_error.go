package capo

import "fmt"

var InternalServerErrorCode = "INTERNAL_ERROR"

// ServerError represents a server error.
type ServerError struct {
	inner error
	Code  string `json:"code"`
}

// NewServerError creates a new instance of server error entity.
func NewServerError(code string, inner error) *ServerError {
	return &ServerError{
		inner: inner,
		Code:  code,
	}
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("%s - %s", e.Code, e.inner.Error())
}
