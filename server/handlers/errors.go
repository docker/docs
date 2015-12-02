package handlers

import (
	"fmt"
)

// VALIDATION ERRORS:

// ErrValidation represents a general validation error
type ErrValidation struct {
	msg string
}

func (err ErrValidation) Error() string {
	return fmt.Sprintf("An error occurred during validation: %s", err.msg)
}

// ErrBadHierarchy represents a missing snapshot at this current time.
// When delegations are implemented it will also represent a missing
// delegation parent
type ErrBadHierarchy struct {
	msg string
}

func (err ErrBadHierarchy) Error() string {
	return fmt.Sprintf("Hierarchy of updates in incorrect: %s", err.msg)
}

// ErrBadRoot represents a failure validating the root
type ErrBadRoot struct {
	msg string
}

func (err ErrBadRoot) Error() string {
	return fmt.Sprintf("The root being updated is invalid: %s", err.msg)
}

// ErrBadTargets represents a failure to validate a targets (incl delegations)
type ErrBadTargets struct {
	msg string
}

func (err ErrBadTargets) Error() string {
	return fmt.Sprintf("The targets being updated is invalid: %s", err.msg)
}

// ErrBadSnapshot represents a failure to validate the snapshot
type ErrBadSnapshot struct {
	msg string
}

func (err ErrBadSnapshot) Error() string {
	return fmt.Sprintf("The snapshot being updated is invalid: %s", err.msg)
}

// END VALIDATION ERRORS
