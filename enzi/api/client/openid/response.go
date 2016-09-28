package openid

import (
	"net/http"

	"github.com/docker/orca/enzi/api/client"
	"github.com/docker/orca/enzi/api/client/openid/oautherrors"
	"github.com/docker/orca/enzi/api/responses"
)

// IDToken is a raw Identity Token string which implements
// (github.com/docker/orca/enzi/api/client).RequestAuthenticator
type IDToken string

var _ client.RequestAuthenticator = IDToken("")

// AuthenticateRequest sets an OpenID JWT identity token authorization header
// on the given request.
func (rawToken IDToken) AuthenticateRequest(r *http.Request) {
	r.Header.Set("Authorization", "Bearer "+string(rawToken))
}

// TokenResponse is the type returned from the OpenID Conect token endpoint. It
// may contain error fields. TokenResponse also implements
// (github.com/docker/orca/enzi/api/client).RequestAuthenticator
type TokenResponse struct {
	*oautherrors.ErrorResponse

	TokenType        string             `json:"token_type,omitempty"         description:"The type of the token; Always 'Bearer'"`
	AccessToken      string             `json:"access_token,omitempty"       description:"Incidentally the same as the identity token"`
	IDToken          string             `json:"id_token,omitempty"           description:"A JWT which can be used to authenticate as the authorized party on behalf of an account to any party specified as an audience which trusts this token issuer"`
	ExpiresIn        int64              `json:"expires_in,omitempty"         description:"The number of seconds from the time of handling this request at which point the token will have expired"`
	RefreshToken     string             `json:"refresh_token,omitempty"      description:"The ID of the subject account. Use the 'refresh_token' grant type to renew an identity token"`
	Account          *responses.Account `json:"account,omitempty"            description:"A description of the subject account"`
	SessionSecret    string             `json:"session_secret,omitempty"     description:"A secret session value which can be used with the 'root_session' or 'service_session' grant type to renew an identity token. When set as a long-lived session token by the service, enables single-sign-on user sessions across services"`
	SessionCSRFToken string             `json:"session_csrf_token,omitempty" description:"A session-specific value against which the service should check a corresponding request cookie to prevent CSRF attacks (if using cookies for session authentication)"`
}

var _ client.RequestAuthenticator = (*TokenResponse)(nil)

// AuthenticateRequest sets an OpenID JWT identity token authorization header
// on the given request.
func (t *TokenResponse) AuthenticateRequest(r *http.Request) {
	IDToken(t.IDToken).AuthenticateRequest(r)
}
