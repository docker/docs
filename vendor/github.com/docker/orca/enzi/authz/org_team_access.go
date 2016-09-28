package authz

import (
	"errors"
	"fmt"

	"github.com/docker/orca/enzi/authn"
	"github.com/docker/orca/enzi/schema"
)

var (
	// ErrAccountAdminAccessRequired indicates that a user does not have
	// admin access for a given account.
	ErrAccountAdminAccessRequired = errors.New("account admin access required")
	// ErrOrgAdminAccessRequired indicates that a user does not have admin
	// access for a given organizaiton.
	ErrOrgAdminAccessRequired = errors.New("organization admin access required")
	// ErrOrgMemberAccessRequired indicates that a user does not have
	// member access for a given organization.
	ErrOrgMemberAccessRequired = errors.New("organization member access required")
	// ErrTeamAdminAccessRequired indicates that a user does not have admin
	// access for a given organizaiton.
	ErrTeamAdminAccessRequired = errors.New("team admin access required")
	// ErrTeamMemberAccessRequired indicates that a user does not have
	// member access for a given organization.
	ErrTeamMemberAccessRequired = errors.New("team member access required")
)

// MembershipAccess holds basic information about a user's permissions in an
// organization or team. This type is ofter used as its pointer type where a
// nil value means the user is not a member.
type MembershipAccess struct {
	IsAdmin bool
}

// AccountAccess determines the given client accounts's permissions in the
// given account. Returns a nil error on success. If the account is an org,
// returns the same result that OrgMembershipAccess would. If the account is a
// user, nil is returned if the given client account is neither a system admin
// or the account in question, otherwise returns a MembershipAccess value with
// IsAdmin set to true.
func (a *authorizer) AccountAccess(acct *schema.Account, clientAccount *authn.Account) (accountAccess *MembershipAccess, err error) {
	// If acct is an org, use OrgMembershipAccess.
	if acct.IsOrg {
		return a.OrgMembershipAccess(acct.ID, clientAccount)
	}

	// Otherwise, the user has Admin access to the account if they are the
	// same user or if the user is a system admin.
	if clientAccount.IsAdmin || clientAccount.ID == acct.ID {
		return &MembershipAccess{IsAdmin: true}, nil
	}

	return nil, nil
}

// OrgMembershipAccess determines the given client account's permissions in the
// given organization. Returns a nil error on success. Returns a nil
// MembershipAccess if the client is not a member of the organizaiton,
// otherwise the client is a member and the value holds whether the client is
// an admin of the organization.
func (a *authorizer) OrgMembershipAccess(orgID string, client *authn.Account) (orgAccess *MembershipAccess, err error) {
	// First, check if the client is a system admin.
	if client.IsAdmin || client.ID == orgID {
		return &MembershipAccess{IsAdmin: true}, nil
	}

	// Next, get their membership in the organizaiton.
	membership, err := a.schemaMgr.GetOrgMembership(orgID, client.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to get org membership: %s", err)
	}

	if membership == nil {
		return nil, nil
	}

	return &MembershipAccess{IsAdmin: membership.IsAdmin}, nil
}

// TeamMembershipAccess determines the given client account's permissions in
// the given team. Returns a nil error on success. Returns a nil
// MembershipAccess if the client is not a member of the team, otherwise the
// client is a member and the value holds whether the client is an admin of the
// team.
func (a *authorizer) TeamMembershipAccess(teamID string, orgAccess MembershipAccess, client *authn.Account) (teamAccess *MembershipAccess, err error) {
	// First, check the client's organization membership. If they are an
	// admin of the org, then they are an admin of the team as well. Note:
	// this also covers the case if they are a system admin.
	if orgAccess.IsAdmin {
		// Admin of the org, automatically admin of the team.
		return &MembershipAccess{IsAdmin: true}, nil
	}

	// Next, get their membership in the team.
	membership, err := a.schemaMgr.GetTeamMembership(teamID, client.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to get team membership: %s", err)
	}

	if membership == nil {
		return nil, nil
	}

	return &MembershipAccess{IsAdmin: membership.IsAdmin}, nil
}
