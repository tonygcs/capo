package capo

import "github.com/tonygcs/gnalog"

var (
	// InternalErrorCode is the error code for the internal errors.
	InternalErrorCode string = "INTERNAL_ERROR"
	// InternalErrorMsg is the error message for the internal errors.
	InternalErrorMsg string = "Internal server error"
	// InvalidRequestBodyCode is error code for bad requests.
	InvalidRequestBodyCode string = "INVALID_REQUEST_BODY"
	// InvalidRequestBodyData is the error message when the request cannot be
	// read by the application.
	InvalidRequestBodyData string = "Cannot read the request body data"
	// InvalidRequestBodyMsg is the error message when the application cannot
	// decode the request body.
	InvalidRequestBodyMsg string = "Invalid request body format"
)

// RequestDecoder is the decoder that the wrapper will use to read each request.
var RequestDecoder = jsonDecoder

// ResponseEncoder is the encoder that the application will use to write the
// information on each response.
var ResponseEncoder = jsonEncoder

// LoggerFactory is the function that will return the logger instance that the
// application will include in the request context.
var LoggerFactory = func() gnalog.Logger {
	return gnalog.New()
}
