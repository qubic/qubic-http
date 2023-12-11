package web

import (
	"github.com/pkg/errors"
)

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// FieldErrorResponse is the form used for API responses when errors
// are encountered during decoding request payload fields.
type FieldErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// FieldsError is used to pass an error during the request through the
// application with web specific context.
type FieldsError struct {
	Err    error
	Fields []FieldError
}

// FieldsError implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (err *FieldsError) Error() string {
	return err.Err.Error()
}

//// NewFieldsError wraps a provided error with an HTTP status code. This
//// function should be used when handlers encounter expected errors.
//func NewFieldsError(err error, status int) error {
//	return &FieldsError{err, status}
//}

// FieldsError is used to pass an error during the request through the
// application with web specific context.
type RequestError struct {
	Err    error
	Status int
}

// NewRequestError wraps a provided error with an HTTP status code. This
// function should be used when handlers encounter expected errors.
func NewRequestError(err error, status int) error {
	return &RequestError{err, status}
}

// FieldsError implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (err *RequestError) Error() string {
	return err.Err.Error()
}

// shutdown is a type used to help with the graceful termination of the service.
type shutdown struct {
	Message string
}

// NewShutdownError returns an error that causes the framework to signal
// a graceful shutdown.
func NewShutdownError(message string) error {
	return &shutdown{message}
}

// FieldsError is the implementation of the error interface.
func (s *shutdown) Error() string {
	return s.Message
}

// IsShutdown checks to see if the shutdown error is contained
// in the specified error value.
func IsShutdown(err error) bool {
	if _, ok := errors.Cause(err).(*shutdown); ok {
		return true
	}
	return false
}
