package enzi

import (
	"fmt"
	"net/http"

	"github.com/docker/orca/auth"
	"github.com/docker/orca/enzi/api"
	"github.com/docker/orca/enzi/api/client"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/forms"
)

func (a *Authenticator) AddTeamMember(ctx *auth.Context, teamID, username string) error {
	session := a.getSession(ctx.ClientCreds)

	if _, err := session.AddTeamMember("id:"+a.defaultOrgID, "id:"+teamID, username, forms.SetMembership{}); err != nil {
		return fmt.Errorf("unable to add team member on auth provider: %s", err)
	}

	return nil
}

func (a *Authenticator) ListTeamMembers(ctx *auth.Context, teamID string) ([]*auth.Account, error) {
	session := a.getSession(ctx.ClientCreds)

	// The interface does not have pagination so ask for thedefault number.
	membersResp, _, err := session.ListTeamMembers("id:"+a.defaultOrgID, "id:"+teamID, "", "", api.MaxPerPageLimit)
	if err != nil {
		apiErrs, ok := err.(*errors.APIErrors)
		if ok && apiErrs.HTTPStatusCode == http.StatusNotFound {
			return nil, auth.ErrTeamDoesNotExist
		}

		return nil, fmt.Errorf("unable to list team members from auth provider: %s", err)
	}

	// NOTE: Do not set team membership on every user.

	accounts := make([]*auth.Account, len(membersResp.Members))
	for i, memberResp := range membersResp.Members {
		accounts[i] = &auth.Account{
			ID:        memberResp.Member.ID,
			FirstName: memberResp.Member.FullName, // FIXME
			LastName:  "",                         // FIXME
			Username:  memberResp.Member.Name,
			Admin:     *memberResp.Member.IsAdmin,
		}
	}

	return accounts, nil
}

func (a *Authenticator) ListUserTeams(ctx *auth.Context, username string) ([]*auth.Team, error) {
	return a.listUserTeams(ctx.ClientCreds, username)
}

func (a *Authenticator) listUserTeams(creds client.RequestAuthenticator, username string) ([]*auth.Team, error) {
	session := a.getSession(creds)

	// The interface does not have pagination so ask for the default
	// number.
	memberTeamsResp, _, err := session.ListOrganizationMemberTeams("id:"+a.defaultOrgID, username, "", api.MaxPerPageLimit)
	if err != nil {
		return nil, fmt.Errorf("unable to get user's teams from auth provider: %s", err)
	}

	teams := make([]*auth.Team, len(memberTeamsResp.MemberTeams))
	for i, memberTeamResp := range memberTeamsResp.MemberTeams {
		// NOTE: We do not get LDAP member sync config for each team
		// because it would be too expensive.
		teams[i] = &auth.Team{
			OrgId:       memberTeamResp.Team.OrgID,
			Name:        memberTeamResp.Team.Name,
			Id:          memberTeamResp.Team.ID,
			Description: memberTeamResp.Team.Description,
		}
	}

	return teams, nil
}

func (a *Authenticator) DeleteTeamMember(ctx *auth.Context, teamID, username string) error {
	session := a.getSession(ctx.ClientCreds)

	if err := session.DeleteTeamMember("id:"+a.defaultOrgID, "id:"+teamID, username); err != nil {
		return fmt.Errorf("unable to delete team member on auth provider: %s", err)
	}

	return nil
}
