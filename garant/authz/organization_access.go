package authz

import (
	"errors"
	"fmt"

	"github.com/docker/dhe-deploy/garant/authn"
)

var (
	// ErrNotAdminOrOrgOwner indicates that a user is not a system admin or an
	// owner of some organization and is therefore not authorized to perform
	// some operation.
	ErrNotAdminOrOrgOwner = errors.New("not a system admin or organization owner")
	// ErrNotAdminOrOrgMember indicates that a user is not a system admin or a
	// member of any teams in some organization and is therefore not
	// authorized to perform some operation.
	ErrNotAdminOrOrgMember = errors.New("not a system admin or organization member")
)

// CheckAdminOrOrgOwner checks whether the given user is either a global admin
// or is a member of the "owners" team of the given organization. Returns nil
// if either is true, ErrNotAdminOrOrgOwner if false, and some other non-nil
// error if there is a problem querying the database.
func (a *authorizer) CheckAdminOrOrgOwner(user *authn.User, orgID string) error {
	if *user.Account.IsAdmin {
		return nil
	}

	// The user is not a system admin, so they must be in the "owners" team of
	// this organization.
	member, err := user.EnziSession.GetOrganizationMember("id:"+orgID, "id:"+user.Account.ID)
	if err != nil {
		return fmt.Errorf("Failed to get details on user's membership in organization: %s", err)
	}

	if member == nil || !member.IsAdmin {
		return ErrNotAdminOrOrgOwner
	}

	return nil
}

// CheckAdminOrOrgMember checks whether the given user is either a global admin
// or is a member of any of the teams in the given organization. Returns nil
// if either is true, ErrNotAdminOrOrgMember if false, and some other non-nil
// error if there is a problem querying the database.
func (a *authorizer) CheckAdminOrOrgMember(user *authn.User, orgID string) error {
	if *user.Account.IsAdmin {
		return nil
	}

	// User is not an admin so ensure that they are a member of the org.
	member, err := user.EnziSession.GetOrganizationMember("id:"+orgID, "id:"+user.Account.ID)
	if err != nil {
		return fmt.Errorf("Failed to get details on user's membership in organization: %s", err)
	}

	if member == nil {
		return ErrNotAdminOrOrgMember
	}

	return nil
}
