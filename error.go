package capo

import (
	"fmt"
)

// capoError is the error data type that the 'Handler' wrapper will use to build
// the error response.
type capoError struct {
	err    error
	status int

	Code string `json:"code"`
	Msg  string `json:"msg"`
}

// Error creates a new error instance with the information that the 'Handler'
// entity is able to understand in order to return an error response to the
// users.
//
// The response to the clients will have the fields 'code' and 'msg'.
func Error(
	err error,
	code string,
	msg string,
	status int,
) *capoError {
	return &capoError{
		err:    err,
		status: status,
		Code:   code,
		Msg:    msg,
	}
}

func (e *capoError) Error() string {
	return fmt.Sprintf(
		"error response with status %d, code %s and msg '%s' :: %v",
		e.status,
		e.Code,
		e.Msg,
		e.err,
	)
}
