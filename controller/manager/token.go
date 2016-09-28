package manager

import (
	"strings"

	"github.com/docker/orca/auth"
)

func (m DefaultManager) AuthenticateSessionToken(tokenStr string) (*auth.Context, error) {
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	return m.GetAuthenticator().AuthenticateSessionToken(tokenStr)
}

func (m DefaultManager) Logout(tokenStr string) error {
	authenticator := m.GetAuthenticator()
	ctx, err := authenticator.AuthenticateSessionToken(tokenStr)
	if err != nil {
		return err
	}

	return authenticator.Logout(ctx)
}
