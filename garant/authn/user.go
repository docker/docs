package authn

import (
	"github.com/docker/garant/auth"
	"github.com/docker/orca/enzi/api/client"
	"github.com/docker/orca/enzi/api/client/openid"
	enziresponses "github.com/docker/orca/enzi/api/responses"
)

const userContextKey = "dtr.authn.User"

// User represents a user account in Docker Trusted Registry.
type User struct {
	// TODO EnziSession and account are both redundant since they can be obtained from the Token
	IsAnonymous bool

	EnziSession *client.Session
	Token       *openid.TokenResponse
	Account     *enziresponses.Account
}

// Assert that *User implements the Garant auth.Account interface.
var _ auth.Account = (*User)(nil)

// Subject returns the account name which will be used as the subject of a JWT.
func (user *User) Subject() string {
	return user.Account.ID
}

// Info returns extra global info about the authenticated account.
func (user *User) Info() map[string]interface{} {
	return map[string]interface{}{
		"name":     user.Account.Name,
		"id":       user.Account.ID,
		"fullName": user.Account.FullName,
		"isAdmin":  *user.Account.IsAdmin,
		"isActive": user.Account.IsActive,
	}
}
