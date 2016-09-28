package errors

import (
	"fmt"
	"net/http"

	"github.com/docker/dhe-deploy/manager/schema"
	"github.com/docker/distribution/context"
)

// APIError describes an API Error.
type APIError struct {
	Code     string      `json:"code"`
	HTTPCode int         `json:"-"`
	Message  string      `json:"message"`
	Detail   interface{} `json:"detail,omitempty"`
}

// Possible error structs
var (
	ErrorCodeInternal = APIError{
		"INTERNAL_ERROR",
		http.StatusInternalServerError,
		"An internal server error occurred. Contact a system administrator for more information.",
		nil,
	}
	ErrorCodeMethodNotAllowed = APIError{
		"METHOD_NOT_ALLOWED",
		http.StatusMethodNotAllowed,
		"The request method is not supported.",
		nil,
	}
	ErrorCodeNotAuthenticated = APIError{
		"NOT_AUTHENTICATED",
		http.StatusUnauthorized,
		"The client is not authenticated.",
		nil,
	}
	ErrorCodeNotAuthorized = APIError{
		"NOT_AUTHORIZED",
		http.StatusForbidden,
		"The client is not authorized.",
		nil,
	}
	ErrorCodeInvalidSettings = APIError{
		"INVALID_SETTINGS",
		http.StatusBadRequest,
		"The submitted settings change request contains invalid values.",
		nil,
	}
	ErrorCodeInvalidAccountName = APIError{
		"INVALID_ACCOUNT_NAME",
		http.StatusBadRequest,
		"The given account name is either too long or contains illegal characters.",
		nil,
	}
	ErrorCodeAccountExists = APIError{
		"ACCOUNT_EXISTS",
		http.StatusBadRequest,
		"An account with the same name already exists.",
		nil,
	}
	ErrorCodeNoSuchAccount = APIError{
		"NO_SUCH_ACCOUNT",
		http.StatusNotFound,
		"An account with the given name does not exist.",
		nil,
	}
	ErrorCodeNoSuchOrganization = APIError{
		"NO_SUCH_ORGANIZATION",
		http.StatusNotFound,
		"An organization with the given name does not exist.",
		nil,
	}
	ErrorCodeNoSuchUser = APIError{
		"NO_SUCH_USER",
		http.StatusNotFound,
		"A user with the given name does not exist.",
		nil,
	}
	ErrorCodeTeamExists = APIError{
		"TEAM_EXISTS",
		http.StatusBadRequest,
		"A team with the same name already exists in the organization.",
		nil,
	}
	ErrorCodeNoSuchTeam = APIError{
		"NO_SUCH_TEAM",
		http.StatusNotFound,
		"A team with the given name does not exist in the organization.",
		nil,
	}
	ErrorCodeNoSuchTag = APIError{
		"NO_SUCH_TAG",
		http.StatusNotFound,
		"A tag with the given name does not exist for the given repository.",
		nil,
	}
	ErrorCodeInvalidJSON = APIError{
		"INVALID_JSON",
		http.StatusBadRequest,
		"Unable to parse JSON",
		nil,
	}
	ErrorCodeInvalidParameters = APIError{
		"INVALID_PARAMETERS",
		http.StatusBadRequest,
		"Unable to parse query parameters",
		nil,
	}
	ErrorCodePasswordIncorrect = APIError{
		"PASSWORD_INCORRECT",
		http.StatusBadRequest,
		"The given password does not match the current password.",
		nil,
	}
	ErrorCodeInvalidRepositoryName = APIError{
		"INVALID_REPOSITORY_NAME",
		http.StatusBadRequest,
		"The given repository name is either too long or contains illegal characters.",
		nil,
	}
	ErrorCodeInvalidRepositoryShortDescription = APIError{
		"INVALID_REPOSITORY_SHORT_DESCRIPTION",
		http.StatusBadRequest,
		"The short description of the repository is invalid.",
		nil,
	}
	ErrorCodeInvalidRepositoryVisibility = APIError{
		"INVALID_REPOSITORY_VISIBILITY",
		http.StatusBadRequest,
		"The visibility value of the repository is invalid.",
		nil,
	}
	ErrorCodeRepositoryExists = APIError{
		"REPOSITORY_EXISTS",
		http.StatusBadRequest,
		"A repository with the same name already exists.",
		nil,
	}
	ErrorCodeReindexJobRunning = APIError{
		"REINDEX_JOB_ALREADY_RUNNING",
		http.StatusBadRequest,
		"The reindex job is already running.",
		nil,
	}
	ErrorCodeNoSuchRepository = APIError{
		"NO_SUCH_REPOSITORY",
		http.StatusNotFound,
		"A repository with the given name does not exist.",
		nil,
	}
	ErrorCodeNoSuchManifest = APIError{
		"NO_SUCH_MANIFEST",
		http.StatusNotFound,
		"A manifest with the given reference does not exist for the given repository.",
		nil,
	}
	ErrorCodeInvalidAccessLevel = APIError{
		"INVALID_ACCESS_LEVEL",
		http.StatusBadRequest,
		"The given access level value is not a valid choice.",
		nil,
	}
	ErrorCodeInvalidRepoContext = APIError{
		"INVALID_REPOSITORY_CONTEXT",
		http.StatusBadRequest,
		"The operation is not valid in the context of the repository.",
		nil,
	}
	ErrorCodeUnableToDeleteTags = APIError{
		"UNABLE_TO_DELETE_TAGS",
		http.StatusInternalServerError,
		"Unable to delete the listed tags",
		nil,
	}
	ErrorCodeTagInNotary = APIError{
		"TAG_IN_NOTARY",
		http.StatusConflict,
		"This tag is in notary and can't be deleted until it is removed from notary",
		nil,
	}
	ErrorCodeOpenID = APIError{
		"OPENID_ERROR",
		http.StatusUnauthorized,
		"Failed to establish openid authentication",
		nil,
	}
	ErrorCodeStatusCheckFailed = APIError{
		"STATUS_CHECK_FAILED",
		http.StatusServiceUnavailable,
		"Failed to check cluster status",
		nil,
	}
)

// InternalError returns an APIError indicating an internal server error.
// Includes the given detail which may be nil. The given error will be logged
// to the logger in the given context, but the resulting APIError will not
// describe that error.
func InternalError(ctx context.Context, err error) APIError {
	context.GetLogger(ctx).Error(err)

	return MakeError(ErrorCodeInternal, map[string]interface{}{
		"requestID": ctx.Value("http.request.id"),
	})
}

func MakeError(template APIError, detail interface{}) APIError {
	return APIError{
		Code:     template.Code,
		HTTPCode: template.HTTPCode,
		Message:  template.Message,
		Detail:   detail,
	}
}

// The following error constructors are used only for errors with special messages in the detail field.
// Use MakeError or the APIError object directly when possible

// NoSuchAccountError returns an APIError indicating that the account with the
// specified name does not exist.
func NoSuchAccountError(name string) APIError {
	return MakeError(ErrorCodeNoSuchAccount, fmt.Sprintf("Account name: %q.", name))
}

func NoSuchOrganizationError(name string) APIError {
	return MakeError(ErrorCodeNoSuchOrganization, fmt.Sprintf("Account name: %q.", name))
}

func NoSuchUserError(name string) APIError {
	return MakeError(ErrorCodeNoSuchUser, fmt.Sprintf("Account name: %q.", name))
}

// NoSuchTeamError returns an APIError indicating that a team with the
// specified name does not exist within the organization.
func NoSuchTeamError(name string) APIError {
	return MakeError(ErrorCodeNoSuchTeam, fmt.Sprintf("Team name: %q.", name))
}

// NoSuchRepositoryError returns an APIError indicating that a repository with
// the specified full name does not exist.
func NoSuchRepositoryError(namespace, name string) APIError {
	fullName := fmt.Sprintf("%s/%s", namespace, name)

	return MakeError(ErrorCodeNoSuchRepository, fmt.Sprintf("Repository name: %q.", fullName))
}

// NoSuchManifestError returns an APIError indicating that a manifest with
// the specified reference does not exist for the specified repository.
func NoSuchManifestError(namespace, name, reference string) APIError {
	fullName := fmt.Sprintf("%s/%s", namespace, name)

	return MakeError(ErrorCodeNoSuchManifest, fmt.Sprintf("Repository name: %q, Reference: %q.", fullName, reference))
}

// NotAuthorizedError returns an APIError indicating that the client is
// not authenticated.
func NotAuthorizedError(detail interface{}) APIError {
	return MakeError(ErrorCodeNotAuthorized, detail)
}

func InvalidSettings(err error) APIError {
	return MakeError(ErrorCodeInvalidSettings, err.Error())
}

// InvalidRepositoryShortDescriptionError returns an APIError indicating that
// a repository short description is invalid.
func InvalidRepositoryShortDescriptionError(maxLength int) APIError {
	return MakeError(ErrorCodeInvalidRepositoryShortDescription, fmt.Sprintf("The maximum length is %d", maxLength))
}

// InvalidRepositoryVisibilityError returns an APIError indicating that a
// repository visibility value is invalid.
func InvalidRepositoryVisibilityError() APIError {
	return MakeError(ErrorCodeInvalidRepositoryVisibility, fmt.Sprintf(
		"Visibiity must be either %q (default) or %q.",
		schema.RepositoryVisibilityPrivate,
		schema.RepositoryVisibilityPublic,
	))
}

// InvalidAccessLevelError returns an APIError indicating that an access level
// was not one of the given valid choices.
func InvalidAccessLevelError(choices ...string) APIError {
	return MakeError(ErrorCodeInvalidAccessLevel, fmt.Sprintf("Access Level must be one of: %q.", choices))
}

// InvalidRepoContextError returns an APIError indicating that an operation is
// not valid in the context of a repository, e.g., an action that may only be
// performed on org-owned repos when it's a user-owned repo.
func InvalidRepoContextError(detail string) APIError {
	return MakeError(ErrorCodeInvalidRepoContext, detail)
}

func UnableToDeleteTagError(tagsWithErrors []string) APIError {
	return MakeError(ErrorCodeUnableToDeleteTags, tagsWithErrors)
}

func OpenIDError(code, description string) APIError {
	return MakeError(ErrorCodeOpenID, fmt.Sprintf("OpenID Connect Error\n\n%s\n\n%s", code, description))
}
