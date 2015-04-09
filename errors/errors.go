package errors

import (
	"fmt"
)

type HTTPError struct {
	HTTPStatus int
	Code       int
	Err        error
}

func (he *HTTPError) Error() string {
	msg := ""
	if he.Err != nil {
		msg = he.Err.Error()
	}
	return fmt.Sprintf("%d: %s", he.Code, msg)
}
