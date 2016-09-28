package util

import (
	"fmt"
	"strings"

	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/integration/apiclient"

	enzierrors "github.com/docker/orca/enzi/api/errors"
	"github.com/stretchr/testify/assert"
)

type ErrorChecker struct {
	errs []string
}

func NewErrorChecker() *ErrorChecker {
	return &ErrorChecker{errs: []string{}}
}

func (ec *ErrorChecker) Errorf(format string, args ...interface{}) {
	ec.errs = append(ec.errs, fmt.Sprintf(format, args...))
}

func (ec *ErrorChecker) Errors() error {
	if len(ec.errs) == 0 {
		return nil
	}
	return fmt.Errorf(strings.Join(ec.errs, "\n"))
}

// RequireErrorCodes requires that the given error is of type
// *apiclient.APIError, that the http status matches the expected HTTP status,
// and that the list of actual errors codes matches the given expected error
// codes.
func (u *Util) RequireErrorCodes(err error, expectedHTTPStatus int, expectedErrors ...errors.APIError) {
	u.RequireErrorCodesWithMsg(err, expectedHTTPStatus, "", expectedErrors...)
}

func (u *Util) RequireErrorCodesWithMsg(err error, expectedHTTPStatus int, msg string, expectedErrors ...errors.APIError) {
	u.AssertErrorCodesWithMsg(err, expectedHTTPStatus, msg, expectedErrors...)
	if u.T().Failed() {
		u.T().FailNow()
	}
}

// AssertErrorCodes is the same as RequireErrorCodes but non-fatal
func (u *Util) AssertErrorCodes(err error, expectedHTTPStatus int, expectedErrors ...errors.APIError) {
	u.AssertErrorCodesWithMsg(err, expectedHTTPStatus, "", expectedErrors...)
}

func (u *Util) AssertErrorCodesWithMsg(err error, expectedHTTPStatus int, msg string, expectedErrors ...errors.APIError) {
	if !assert.NotNil(u.T(), err, msg) {
		return
	}
	if !assert.IsType(u.T(), &apiclient.APIError{}, err, msg) {
		assert.Fail(u.T(), "found error of unexpected type", "actual: %s\nmessage: %s", err.Error(), msg)
		return
	}

	errWrapper := err.(*apiclient.APIError)
	assert.Equal(u.T(), expectedHTTPStatus, errWrapper.HTTPStatusCode, fmt.Sprintf("unexpected HTTP status code; %s", msg))

	// Build a set of actual and expected error codes.
	actualErrorCodes := make(map[string]struct{}, len(errWrapper.Errors))
	for _, err := range errWrapper.Errors {
		actualErrorCodes[err.Code] = struct{}{}
	}

	expectedErrorCodes := make(map[string]struct{}, len(expectedErrors))
	for _, err := range expectedErrors {
		expectedErrorCodes[err.Code] = struct{}{}
	}

	missingOrUnexpected := false

	// Loop through the actual error code set. If a code exists in the
	// expected code set then delete it from the expected error code set.
	// If it is nou.T, fail and log the unexpected error code.
	for code := range actualErrorCodes {
		if _, ok := expectedErrorCodes[code]; ok {
			delete(expectedErrorCodes, code)
		} else {
			missingOrUnexpected = true
		}
	}

	// Any remaning codes in the expected error code set are those which
	// should have occurred but did not.
	if len(expectedErrorCodes) > 0 {
		missingOrUnexpected = true
	}

	if missingOrUnexpected {
		assert.Fail(u.T(), "error codes mismatch", "expected: %v\n actual: %v\nmessage: %s", expectedErrors, err, msg)
	}
}

func (u *Util) RequireEnziErrorCode(err error, expectedHTTPStatus int, expectedError *enzierrors.APIError) {
	u.RequireEnziErrorCodeWithMsg(err, expectedHTTPStatus, "", expectedError)
}

func (u *Util) RequireEnziErrorCodeWithMsg(err error, expectedHTTPStatus int, msg string, expectedError *enzierrors.APIError) {
	u.AssertEnziErrorCodeWithMsg(err, expectedHTTPStatus, msg, expectedError)
	if u.T().Failed() {
		u.T().FailNow()
	}
}

// AssertErrorCodes is the same as RequireErrorCodes but non-fatal
func (u *Util) AssertEnziErrorCode(err error, expectedHTTPStatus int, expectedError *enzierrors.APIError) {
	u.AssertEnziErrorCodeWithMsg(err, expectedHTTPStatus, "", expectedError)
}

func (u *Util) AssertEnziErrorCodeWithMsg(err error, expectedHTTPStatus int, msg string, expectedError *enzierrors.APIError) {
	if !assert.NotNil(u.T(), err, msg) {
		return
	}
	if !assert.IsType(u.T(), &enzierrors.APIErrors{}, err, msg) {
		assert.Fail(u.T(), "found error of unexpected type", "actual: %s\nmessage: %s", err.Error(), msg)
		return
	}

	apiErrs := err.(*enzierrors.APIErrors)
	assert.Equal(u.T(), expectedHTTPStatus, apiErrs.HTTPStatusCode, fmt.Sprintf("unexpected HTTP status code; %s", msg))

	for _, apiErr := range apiErrs.Errors {
		assert.Equal(u.T(), apiErr.Code, expectedError.Code)
		if msg != "" {
			assert.Equal(u.T(), apiErr.Message, expectedError.Message)
		}
	}
}
