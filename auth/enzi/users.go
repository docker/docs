package enzi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/orca/auth"
	"github.com/docker/orca/enzi/api"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/forms"
)

func (a *Authenticator) SaveUser(ctx *auth.Context, user *auth.Account) (string, error) {
	// Check if the user exists. If it exists then do an update. If it
	// does not exist then create it.
	_, err := a.GetUser(ctx, user.Username)
	if err == nil {
		return "update-account", a.updateUser(ctx, user)
	}
	if err == auth.ErrAccountDoesNotExist {
		return "add-account", a.createUser(ctx, user)
	}

	// Pass the error through.
	return "", err
}

func (a *Authenticator) updateUser(ctx *auth.Context, user *auth.Account) error {
	session := a.getSession(ctx.ClientCreds)

	// eNZi only has a full name field for users, not first and last name
	// fields.
	fullName := strings.TrimSpace(fmt.Sprintf("%s %s", user.FirstName, user.LastName))

	updateAccountForm := forms.UpdateAccount{
		FullName: &fullName,
	}

	if ctx.User.Admin {
		updateAccountForm.IsAdmin = &user.Admin
		isActive := !user.Disabled
		updateAccountForm.IsActive = &isActive
	}

	if _, err := session.UpdateAccount(user.Username, updateAccountForm); err != nil {
		return fmt.Errorf("unable to update user on auth provider: %s", err)
	}

	// Update the user's default role ONLY if the authenticated user is
	// an admin.
	if ctx.User.Admin {
		if err := a.SetUserRole(user.Username, user.Role); err != nil {
			return fmt.Errorf("unable to set user's default role: %s", err)
		}
	}

	// Handle create, update (labels), and delete of public keys.
	if err := a.updateUserPublicKeys(user); err != nil {
		return fmt.Errorf("unable to update user's public keys: %s", err)
	}

	// Handle an admin changing another user's password.
	if ctx.User.Admin && ctx.User.Username != user.Username && user.Password != "" {
		// The old password is not required if an admin is changing
		// another user's password.
		if err := a.ChangePassword(ctx, user.Username, "", user.Password); err != nil {
			return fmt.Errorf("unable to change user's password: %s", err)
		}
	}

	return nil
}

func (a *Authenticator) createUser(ctx *auth.Context, user *auth.Account) error {
	session := a.getSession(ctx.ClientCreds)

	// eNZi only has a full name field for users, not first and last name
	// fields.
	fullName := strings.TrimSpace(fmt.Sprintf("%s %s", user.FirstName, user.LastName))

	createAccountForm := forms.CreateAccount{
		Name:     user.Username,
		FullName: fullName,
		IsOrg:    false,
		IsAdmin:  user.Admin,
		IsActive: !user.Disabled,
		// The auth service will handle hashing the password.
		Password: user.Password,
	}

	userResp, err := session.CreateAccount(createAccountForm)
	if err != nil {
		return fmt.Errorf("unable to create user on auth provider: %s", err)
	}

	// Set the user ID from the response.
	user.ID = userResp.ID

	// Set the user's default role ONLY if the authenticated user is
	// an admin.
	if ctx.User.Admin {
		if err := a.SetUserRole(user.Username, user.Role); err != nil {
			return fmt.Errorf("unable to set user's default role: %s", err)
		}
	}

	return nil
}

func (a *Authenticator) GetUser(ctx *auth.Context, username string) (*auth.Account, error) {
	session := a.getSession(ctx.ClientCreds)

	accountResp, err := session.GetAccount(username)
	if err != nil {
		apiErrs, ok := err.(*errors.APIErrors)
		if ok && apiErrs.HTTPStatusCode == http.StatusNotFound {
			return nil, auth.ErrAccountDoesNotExist
		}

		return nil, fmt.Errorf("unable to get user from auth provider: %s", err)
	}

	if accountResp.IsOrg {
		// This interface doesn't have Organizations.
		return nil, auth.ErrAccountDoesNotExist
	}

	return a.populateUserFields(ctx.ClientCreds, accountResp)
}

func (a *Authenticator) ListUsers(ctx *auth.Context) ([]*auth.Account, error) {
	session := a.getSession(ctx.ClientCreds)

	// Filter to users only because this interface doesn't have orgs. It
	// also doesn't support pagination so ask for the default number.
	accountsResp, _, err := session.ListAccounts("users", "", api.MaxPerPageLimit)
	if err != nil {
		return nil, fmt.Errorf("unable to get users from auth provider: %s", err)
	}

	// NOTE: Do not set team membership on every user.

	accounts := make([]*auth.Account, len(accountsResp.Accounts))
	for i, accountResp := range accountsResp.Accounts {
		accounts[i] = &auth.Account{
			ID:        accountResp.ID,
			FirstName: accountResp.FullName, // FIXME
			LastName:  "",                   // FIXME
			Username:  accountResp.Name,
			Admin:     *accountResp.IsAdmin,
		}
	}

	return accounts, nil
}

func (a *Authenticator) DeleteUser(ctx *auth.Context, user *auth.Account) error {
	session := a.getSession(ctx.ClientCreds)

	if err := session.DeleteAccount(user.Username); err != nil {
		return fmt.Errorf("unable to delete user on auth provider: %s", err)
	}

	return nil
}
