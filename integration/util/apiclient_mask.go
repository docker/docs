package util

import (
	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/api/responses"
)

func (u *Util) CreateUser(username, password string) (*responses.Account, error) {
	if u.IsSuiteRunningInLDAPMode() {
		err := u.CreateUserInLDAPServer(username, password)
		if err != nil {
			return nil, err
		}

		user, err := u.API.EnziSession().CreateAccount(forms.CreateAccount{
			Name:     username,
			Password: password,
		})
		if err != nil {
			u.DeleteLDAPUser(username)
		}
		return user, err
	}

	return u.API.EnziSession().CreateAccount(forms.CreateAccount{
		Name:     username,
		Password: password,
	})
}

func (u *Util) ActivateUser(username string) error {
	if u.IsSuiteRunningInManagedMode() {
		_, err := u.API.EnziSession().UpdateAccount(username, forms.UpdateAccount{
			IsActive: &[]bool{true}[0],
		})
		return err
	}

	// Activating a user in LDAP mode is not supported
	return nil
}

func (u *Util) DeactivateUser(username string) error {
	if u.IsSuiteRunningInManagedMode() {
		_, err := u.API.EnziSession().UpdateAccount(username, forms.UpdateAccount{
			IsActive: &[]bool{false}[0],
		})
		return err
	}

	// Deactivating a user in LDAP mode is not supported
	return nil
}

func (u *Util) DeleteAccount(username string) error {
	if u.IsSuiteRunningInLDAPMode() {
		u.DeleteLDAPEntry(username)
	}

	return u.API.EnziSession().DeleteAccount(username)
}

func (u *Util) ChangePassword(username, oldPassword, newPassword string) error {
	if u.IsSuiteRunningInLDAPMode() {
		return u.ChangeUserLDAPPassword(username, oldPassword, newPassword)
	}

	_, err := u.API.EnziSession().ChangePassword(username, forms.ChangePassword{oldPassword, newPassword})
	return err
}
