package tmperrors

// This package is temporary. It should be merged with enzi's responses package when the jobrunner is moved back into enzi

import (
	"fmt"
	"net/http"

	"github.com/docker/orca/enzi/api/errors"
)

// NoSuchCron returns an error indicating that a cron with the given action does
// not exist.
func NoSuchCron(action string) *errors.APIError {
	return &errors.APIError{
		Code:     "NO_SUCH_CRON",
		HTTPCode: http.StatusNotFound,
		Message:  "A cron with the given action does not exist.",
		Detail:   fmt.Sprintf("Cron action: %q", action),
	}
}

// NoSuchActionConfig returns an error indicating that a actionConfig with the given action does
// not exist.
func NoSuchActionConfig(action string) *errors.APIError {
	return &errors.APIError{
		Code:     "NO_SUCH_ACTION_CONFIG",
		HTTPCode: http.StatusNotFound,
		Message:  "An action config with the given action does not exist.",
		Detail:   fmt.Sprintf("ActionConfig action: %q", action),
	}
}
