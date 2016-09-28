package oautherrors

// Response includes the fields of a standard Oauth error response: an error
// code and description.
type Response struct {
	ErrorCode   string `json:"error,omitempty"             description:"The error code"`
	Description string `json:"error_description,omitempty" description:"The error description"`
}

func (r Response) Error() string {
	return r.Description
}

// Standard Oauth Error Codes
const (
	// CodeInvalidRequest indicates that the request is missing a required
	// parameter, includes an unsupported parameter value (other than grant
	// type), repeats a parameter, includes multiple credentials, utilizes
	// more than one mechanism for authenticating the client, or is
	// otherwise malformed.
	CodeInvalidRequest = "invalid_request"
	// CodeInvalidClient indicates that client authentication failed (e.g.,
	// unknown client, no client authentication included, or unsupported
	// authentication method).  The authorization server MAY return an HTTP
	// 401 (Unauthorized) status code to indicate which HTTP authentication
	// schemes are supported.  If the client attempted to authenticate via
	// the "Authorization" request header field, the authorization server
	// MUST respond with an HTTP 401 (Unauthorized) status code and
	// include the "WWW-Authenticate" response header field matching the
	// authentication scheme used by the client.
	CodeInvalidClient = "invalid_client"
	// CodeInvalidGrant indicates that the provided authorization grant
	// (e.g., authorization code, resource owner credentials) or refresh
	// token is invalid, expired, revoked, does not match the redirection
	// URI used in the authorization request, or was issued to another
	// client.
	CodeInvalidGrant = "invalid_grant"
	// CodeUnauthorizedClient indicates that the authenticated client is
	// not authorized to use this authorization grant type.
	CodeUnauthorizedClient = "unauthorized_client"
	// CodeUnsupportedGrantType indicates that the authorization grant type
	// is not supported by the authorization server.
	CodeUnsupportedGrantType = "unsupported_grant_type"
)

func oauthError(errorCode string, err error) Response {
	return Response{
		ErrorCode:   errorCode,
		Description: err.Error(),
	}
}

// InvalidRequest returns an Oauth Error Response indicating that the request
// is invalid. The given error will be used as the description.
func InvalidRequest(err error) Response {
	return oauthError(CodeInvalidRequest, err)
}

// InvalidClient returns an Oauth Error Response indicating that client
// authentication failed. The given error will be used as the description.
func InvalidClient(err error) Response {
	return oauthError(CodeInvalidClient, err)
}

// InvalidGrant returns an Oauth Error Response indicating that the grant type
// used is not valid for the client. The given error will be used as the
// description.
func InvalidGrant(err error) Response {
	return oauthError(CodeInvalidGrant, err)
}

// UnauthorizedClient returns an Oauth Error Response indicating that the
// authenticated client is not authorized to use the requested grant type. The
// given error will be used as the description.
func UnauthorizedClient(err error) Response {
	return oauthError(CodeUnauthorizedClient, err)
}

// UnsupportedGrantType returns an Oauth Error Response indicating that the
// authorization grant type is not supported. The given error will be used as
// the description.
func UnsupportedGrantType(err error) Response {
	return oauthError(CodeUnsupportedGrantType, err)
}
