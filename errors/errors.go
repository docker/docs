package errors

import (
	"fmt"
)

// HTTPError represents an application error which will map to
// an HTTP status code and returned error object.
type HTTPError struct {
	HTTPStatus int
	Code       int
	Err        error
}

// Error implements the error interface
func (he *HTTPError) Error() string {
	msg := ""
	if he.Err != nil {
		msg = he.Err.Error()
	}
	return fmt.Sprintf("%d: %s", he.Code, msg)
}
