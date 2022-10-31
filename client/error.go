package client

import "fmt"

// ServerError represents an error form the server side.
type ServerError struct {
	status int
	data   []byte
}

// newServerError creates a new instance of server error.
func newServerError(status int, data []byte) *ServerError {
	return &ServerError{
		status: status,
		data:   data,
	}
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("%d", e.status)
}

// Read reads the server error response and set the data in the entity provided.
func (e *ServerError) Read(entity interface{}) error {
	return m.Unmarshal(e.data, entity)
}
