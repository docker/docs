package enzi

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/config"
)

func (a *Authenticator) CanChangePassword(ctx *auth.Context) bool {
	session := a.getSession(ctx.ClientCreds)

	authConfigResp, err := session.GetAuthConfig()
	if err != nil {
		log.Errorf("unable to get auth config from provider: %s", err)
		return false
	}

	return authConfigResp.Backend == config.AuthBackendManaged
}

func (a *Authenticator) ChangePassword(ctx *auth.Context, username, oldPassword, newPassword string) error {
	session := a.getSession(ctx.ClientCreds)

	changePasswordForm := forms.ChangePassword{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	if _, err := session.ChangePassword(username, changePasswordForm); err != nil {
		return fmt.Errorf("unable to change user's password on auth provider: %s", err)
	}

	return nil
}
