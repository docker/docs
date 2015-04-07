package errors

import (
	"fmt"
)

type DockerError struct {
	HTTPStatus int
	Code       int
	Error      error
}

func (de *DockerError) Error() string {
	fmt.Sprintf("%d: %s", de.Code, de.Error.Error())
}
