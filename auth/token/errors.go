package token

import (
	"errors"
)

// Errors used and exported by this package.
var (
	ErrInsufficientScope = errors.New("insufficient scope")
	ErrTokenRequired     = errors.New("authorization token required")
	ErrMalformedToken    = errors.New("malformed token")
	ErrInvalidToken      = errors.New("invalid token")
)
