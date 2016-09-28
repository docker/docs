package manager

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/auth"
)

var (
	ksAccounts  = datastoreVersion + "/accounts"
	ksTokens    = datastoreVersion + "/tokens"
	ksBlacklist = datastoreVersion + "/blacklist"
)

func (m DefaultManager) Account(ctx *auth.Context, username string) (*auth.Account, error) {
	return m.GetAuthenticator().GetUser(ctx, username)
}

func (m DefaultManager) Accounts(ctx *auth.Context) ([]*auth.Account, error) {
	return m.GetAuthenticator().ListUsers(ctx)
}

func (m DefaultManager) AccountsForTeam(ctx *auth.Context, teamID string) ([]*auth.Account, error) {
	return m.GetAuthenticator().ListTeamMembers(ctx, teamID)
}

func (m DefaultManager) SaveAccount(ctx *auth.Context, updateAccount *auth.Account) (string, error) {
	requester := ctx.User

	// Do some sanity checking of the request to make sure we should let it through
	if !requester.Admin {
		// Non admin's can't modify someone else
		if requester.Username != updateAccount.Username {
			// XXX Shoudl we send a more formal event?
			log.Warnf("Non-admin user %s attempted to modify a different user %s", requester.Username, updateAccount.Username)
			return "", auth.ErrUnauthorized
		}

		// If we get here, this is a self modification.  We allow
		// PublicKey modifications, but everything else must be
		// identical.  We could explicity try to check fields one by
		// one and error invalid requests, but that increases the
		// chance of regressions slipping through.  Instead we'll just
		// clobber the user input except for the field(s) we allow
		// changing, and we'll silently discard changes that aren't
		// permitted.  without erroring.
		requester.PublicKeys = updateAccount.PublicKeys
		updateAccount = requester
		// Make sure not to clobber password
		updateAccount.Password = ""
	}

	eventType, err := m.GetAuthenticator().SaveUser(ctx, updateAccount)
	if err != nil {
		return "", err
	}
	m.logEvent(eventType, fmt.Sprintf("username=%s", updateAccount.Username), []string{"security"})
	return eventType, nil
}

func (m DefaultManager) DeleteAccount(ctx *auth.Context, account *auth.Account) error {
	if err := m.GetAuthenticator().DeleteUser(ctx, account); err != nil {
		return err
	}
	m.logEvent("delete-account", fmt.Sprintf("username=%s", account.Username), []string{"security"})
	return nil
}

func (m DefaultManager) ChangePassword(ctx *auth.Context, username, oldPassword, newPassword string) error {
	authenticator := m.GetAuthenticator()
	if !authenticator.CanChangePassword(ctx) {
		return fmt.Errorf("not supported for authenticator: %s", authenticator.Name())
	}

	if err := authenticator.ChangePassword(ctx, username, oldPassword, newPassword); err != nil {
		return err
	}
	m.logEvent("change-password", username, []string{"security"})
	return nil
}
