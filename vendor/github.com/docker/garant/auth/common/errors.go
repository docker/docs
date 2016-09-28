package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
)

// ClientError exists to differentiate from server errors.
type ClientError struct {
	Err  error
	Code int
}

// Errors exported by this package.
var (
	ErrAuthenticationRequired    = &ClientError{Err: errors.New("authentication required"), Code: http.StatusUnauthorized}
	ErrMalformedAuthHeader       = &ClientError{Err: errors.New("malformed HTTP Authorization header")}
	ErrIncorrectUsernamePassword = &ClientError{Err: errors.New("incorrect username or password"), Code: http.StatusUnauthorized}
	ErrInvalidToken              = &ClientError{Err: errors.New("invalid token"), Code: http.StatusUnauthorized}
	ErrInvalidGrantType          = &ClientError{Err: errors.New("invalid grant type"), Code: http.StatusUnauthorized}
	ErrInvalidClientID           = &ClientError{Err: errors.New("invalid client ID")}
	ErrTooManyFailedLogins       = &ClientError{Err: errors.New("too many failed login attempts for username or IP address"), Code: 429}
)

// StatusCode returns the response status for this client error.
func (ce *ClientError) StatusCode() int {
	if ce.Code == 0 {
		return 400
	}

	return ce.Code
}

func (ce *ClientError) Error() string {
	return ce.Err.Error()
}

func (ce *ClientError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ce.StatusCode() == http.StatusUnauthorized {
		w.Header().Add("WWW-Authenticate", fmt.Sprintf("Basic realm=%q", r.Host))
	}

	w.WriteHeader(ce.StatusCode())

	encoder := json.NewEncoder(w)
	encoder.Encode(map[string]string{"details": ce.Error()})
}

// Stacker is any object which can return a stack trace.
type Stacker interface {
	Stack() []byte
}

type stackTraceError struct {
	err   error
	stack []byte
}

// WithStackTrace wraps the given error and attaches a stack trace of the
// calling goroutine.
func WithStackTrace(err error) error {
	if _, ok := err.(Stacker); ok {
		return err // Already has a stack trace.
	}

	return &stackTraceError{
		err: err,
		// TODO: `debug.Stack()` is deprecated in favor of `runtime.Stack(...)`
		// but it's easier to use for now (no need to give a buffer).
		stack: debug.Stack(),
	}
}

func (se *stackTraceError) Error() string {
	return se.err.Error()
}

func (se *stackTraceError) Stack() []byte {
	return se.stack
}
