package authn

import (
	"github.com/docker/orca/enzi/schema"
)

// Account represents an authenticated account in eNZi.
type Account struct {
	schema.Account
	IsAnonymous bool

	// When authenticated via a session cookie.
	Session *schema.Session

	// When authenticated via a token issued to a service.
	AuthorizedService *schema.Service
}

// MakeAnonymousAccount creates an anonymous account to indicate that the
// client is not authenticated.
func MakeAnonymousAccount() *Account {
	return &Account{
		Account: schema.Account{
			Name: "anonymous",
		},
		IsAnonymous: true,
	}
}
