package errors

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/docker/distribution/context"
)

// APIError describes an API Error.
type APIError struct {
	Code     string      `json:"code"`
	HTTPCode int         `json:"-"`
	Message  string      `json:"message"`
	Detail   interface{} `json:"detail,omitempty"`
}

func (err *APIError) Error() string {
	errMsg := fmt.Sprintf("%s: %s", err.Code, err.Message)
	if err.Detail != nil {
		errMsg = fmt.Sprintf("%s - Detail: %v", errMsg, err.Detail)
	}

	return errMsg
}

// IsInternal returns whether this APIError is an internal error.
func (err *APIError) IsInternal() bool {
	return err.Code == "INTERNAL_ERROR"
}

// APIErrors is a response object which wraps one or more API errors.
type APIErrors struct {
	HTTPStatusCode int         `json:"-"`
	Errors         []*APIError `json:"errors"`
}

func (errs *APIErrors) Error() string {
	messages := make([]string, len(errs.Errors))
	for i, err := range errs.Errors {
		messages[i] = err.Error()
	}

	return fmt.Sprintf("API Errors: %q", messages)
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

// Internal returns an APIError indicating an internal server error. Includes
// the given detail which may be nil. The given error will be logged to the
// logger in the given context, but the resulting APIError will not describe
// that error.
func Internal(ctx context.Context, err error) *APIError {
	err = WithStackTrace(err)
	ctx = context.WithValue(ctx, "stackTrace", string(err.(Stacker).Stack()))
	context.GetLogger(ctx, "stackTrace").Error(err.Error())

	return &APIError{
		Code:     "INTERNAL_ERROR",
		HTTPCode: http.StatusInternalServerError,
		Message:  "An internal server error occurred. Contact a system administrator for more information.",
		Detail: map[string]interface{}{
			"requestID": ctx.Value("http.request.id"),
		},
	}
}

// HandlerNotFound returns an error indicating that the requested HTTP method
// and path have no registered handler.
func HandlerNotFound() *APIError {
	return &APIError{
		Code:     "HANDLER_NOT_FOUND",
		HTTPCode: http.StatusNotFound,
		Message:  "No registered handler for this method and path.",
	}
}

// AuthenticationRequired returns an error indicating that the requested API
// action requires authentication which was not supplied or not successful. The
// given reason string will be used as the detail for the error.
func AuthenticationRequired() *APIError {
	return &APIError{
		Code:     "AUTHENTICATION_REQUIRED",
		HTTPCode: http.StatusUnauthorized,
		Message:  "The client is required to authenticate.",
	}
}

// InactiveAccount returns an error indicating that the client could not
// authenticate because their account is currentnly inactive.
func InactiveAccount() *APIError {
	return &APIError{
		Code:     "INACTIVE_ACCOUNT",
		HTTPCode: http.StatusUnauthorized,
		Message:  "Your account is currently inactive.",
	}
}

// NotAuthorized returns an error indicating that the client is not authorized
// to perform the requested API action. The given reason may explain what is
// required for the client to be authorized.
func NotAuthorized(reason string) *APIError {
	return &APIError{
		Code:     "NOT_AUTHORIZED",
		HTTPCode: http.StatusForbidden,
		Message:  "The client is not authorized.",
		Detail:   reason,
	}
}

// InvalidJSON returns an error indicating that the client has submitted a JSON
// form that the server was unable to parse.
func InvalidJSON(jsonDecodeErr error) *APIError {
	return &APIError{
		Code:     "INVALID_JSON",
		HTTPCode: http.StatusBadRequest,
		Message:  "Unable to parse JSON",
		Detail:   jsonDecodeErr.Error(),
	}
}

// InvalidFormField returns an error indicating that the client has submitted
// an invalid form field.
func InvalidFormField(fieldName, reason string) *APIError {
	return &APIError{
		Code:     "INVALID_FORM_FIELD",
		HTTPCode: http.StatusBadRequest,
		Message:  "An input form field is invalid",
		Detail: map[string]string{
			"field":  fieldName,
			"reason": reason,
		},
	}
}

// InvalidAuthentication returns an error indicating that the client has
// submitted invalid authentication credentials.
func InvalidAuthentication(detail string) *APIError {
	return &APIError{
		Code:     "INVALID_AUTHENTICATION_CREDENTIALS",
		HTTPCode: http.StatusUnauthorized,
		Message:  "invalid authentication credentials given",
		Detail:   detail,
	}
}

// AccountExists returns an error indicating that an account already exists
// with the same name as one which the client is attempting to create.
func AccountExists() *APIError {
	return &APIError{
		Code:     "ACCOUNT_EXISTS",
		HTTPCode: http.StatusBadRequest,
		Message:  "An account with the same name already exists.",
	}
}

// AccountsExists returns an error indicating that one or more of the names of
// the accounts that the user is attempting to create already exist.
func AccountsExist(dupNames []string) *APIError {
	return &APIError{
		Code:     "ACCOUNTS_EXIST",
		HTTPCode: http.StatusBadRequest,
		Message:  "Accounts with the same names already exist.",
		Detail:   map[string]interface{}{"duplicateNames": dupNames},
	}
}

// NoSuchAccount returns an error indicating that an account with the given
// name or ID does not exist.
func NoSuchAccount(nameOrID string) *APIError {
	return &APIError{
		Code:     "NO_SUCH_ACCOUNT",
		HTTPCode: http.StatusNotFound,
		Message:  "An account with the given name or ID does not exist.",
		Detail:   fmt.Sprintf("Account name or ID: %q", nameOrID),
	}
}

// CannotChangePassword returns an error indicating that the client cannot
// change a password for the given reason.
func CannotChangePassword(reason string) *APIError {
	return &APIError{
		Code:     "CANNOT_CHANGE_PASSWORD",
		HTTPCode: http.StatusBadRequest,
		Message:  "Unable to change account password.",
		Detail:   reason,
	}
}

// PasswordIncorrect returns an error indicating that a given password does not
// match the current password for a user.
func PasswordIncorrect() *APIError {
	return &APIError{
		Code:     "PASSWORD_INCORRECT",
		HTTPCode: http.StatusBadRequest,
		Message:  "The given password does not match the current password.",
	}
}

// CannotCreateUser returns an error indicating that a user account cannot be
// created. The given reason is used as the detail for the error.
func CannotCreateUser(reason string) *APIError {
	return &APIError{
		Code:     "CANNOT_CREATE_USER",
		HTTPCode: http.StatusBadRequest,
		Message:  "Unable create user account.",
		Detail:   reason,
	}
}

// NoSuchMember returns an error indicating that a user with the given
// name or ID is not a member of some org or team.
func NoSuchMember(nameOrID string) *APIError {
	return &APIError{
		Code:     "NO_SUCH_MEMBER",
		HTTPCode: http.StatusNotFound,
		Message:  "An user with the given name or ID is not a member of the organization or team.",
		Detail:   fmt.Sprintf("User name or ID: %q", nameOrID),
	}
}

// TeamExists returns an error indicating that a team with the given name
// already exists when the client is attempting to create a new team.
func TeamExists() *APIError {
	return &APIError{
		Code:     "TEAM_EXISTS",
		HTTPCode: http.StatusBadRequest,
		Message:  "A team with the same name already exists in the organization.",
	}
}

// NoSuchTeam returns an error indicating that a team with the given name or ID
// does not exist.
func NoSuchTeam(nameOrID string) *APIError {
	return &APIError{
		Code:     "NO_SUCH_TEAM",
		HTTPCode: http.StatusNotFound,
		Message:  "A team with the given name or ID does not exist in the organization.",
		Detail:   fmt.Sprintf("Team name or ID: %q", nameOrID),
	}
}

// LdapPrecludes returns an error indicating that the client cannot complete an
// operation because current LDAP configuration precludes that ability. The
// given reason string is used as detail in the erorr.
func LdapPrecludes(reason string) *APIError {
	return &APIError{
		Code:     "LDAP_PRECLUDES",
		HTTPCode: http.StatusConflict,
		Message:  "LDAP configuration prevents this action",
		Detail:   reason,
	}
}

// ServiceExists returns an error indicating that a service already exists
// with the same name as one which the client is attempting to create.
func ServiceExists() *APIError {
	return &APIError{
		Code:     "SERVICE_EXISTS",
		HTTPCode: http.StatusBadRequest,
		Message:  "A service with the same name already exists.",
	}
}

// NoSuchService returns an error indicating that a service with the given
// name or ID does not exist.
func NoSuchService(nameOrID string) *APIError {
	return &APIError{
		Code:     "NO_SUCH_SERVICE",
		HTTPCode: http.StatusNotFound,
		Message:  "A service with the given name or ID does not exist.",
		Detail:   fmt.Sprintf("Service name or ID: %q", nameOrID),
	}
}

// NoSuchJob returns an error indicating that a job with the given jobID does
// not exist.
func NoSuchJob(jobID string) *APIError {
	return &APIError{
		Code:     "NO_SUCH_JOB",
		HTTPCode: http.StatusNotFound,
		Message:  "A job with the given ID does not exist.",
		Detail:   fmt.Sprintf("Job ID: %q", jobID),
	}
}

// NoSuchWorker returns an error indicating that a worker with the given
// workerID does not exist.
func NoSuchWorker(workerID string) *APIError {
	return &APIError{
		Code:     "NO_SUCH_WORKER",
		HTTPCode: http.StatusNotFound,
		Message:  "A worker with the given ID does not exist.",
		Detail:   fmt.Sprintf("Worker ID: %q", workerID),
	}
}

// WorkerUnavailable returns an error indicating that API server was unable to
// complete a request to the worker with the given ID due to the given error.
func WorkerUnavailable(ctx context.Context, workerID string, err error) *APIError {
	context.GetLogger(ctx).Warnf("unable to perform worker request: %s", err)

	return &APIError{
		Code:     "WORKER_UNAVAILABLE",
		HTTPCode: http.StatusServiceUnavailable,
		Message:  "Unable to make request to worker.",
		Detail: map[string]interface{}{
			"workerID": workerID,
			"error":    err.Error(),
		},
	}
}
