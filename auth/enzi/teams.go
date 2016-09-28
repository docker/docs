package enzi

import (
	"fmt"
	"net/http"

	"github.com/docker/orca/auth"
	"github.com/docker/orca/enzi/api"
	"github.com/docker/orca/enzi/api/client"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/api/responses"
)

func (a *Authenticator) SaveTeam(ctx *auth.Context, team *auth.Team) (string, error) {
	session := a.getSession(ctx.ClientCreds)

	// Check if the team exists. If it exists then do an update. If it
	// does not exist then create it.
	// FIXME: Is this correct? The fact that it even has an ID ... wouldn't
	// that mean that it already exists?
	_, err := a.GetTeam(ctx, team.Id)
	if err == nil {
		return "update-team", updateTeam(session, a.defaultOrgID, team)
	}
	if err == auth.ErrTeamDoesNotExist {
		return "add-team", createTeam(session, a.defaultOrgID, team)
	}

	// Pass the error through.
	return "", err
}

func updateTeam(session *client.Session, orgID string, team *auth.Team) error {
	updateTeamForm := forms.UpdateTeam{
		Name:        &team.Name,
		Description: &team.Description,
	}

	if _, err := session.UpdateTeam("id:"+orgID, "id:"+team.Id, updateTeamForm); err != nil {
		return fmt.Errorf("unable to update team on auth provider: %s", err)
	}

	// Handle LDAP Sync Options.
	if err := setTeamMemberSyncConfig(session, orgID, team); err != nil {
		return err
	}

	return nil
}

func createTeam(session *client.Session, orgID string, team *auth.Team) error {
	createTeamForm := forms.CreateTeam{
		Name:        team.Name,
		Description: team.Description,
	}

	teamResp, err := session.CreateTeam("id:"+orgID, createTeamForm)
	if err != nil {
		return fmt.Errorf("unable to create team on auth provider: %s", err)
	}

	// Set the team ID from the response.
	team.Id = teamResp.ID

	// Handle LDAP Sync Options.
	if err := setTeamMemberSyncConfig(session, orgID, team); err != nil {
		return err
	}

	return nil
}

func setTeamMemberSyncConfig(session *client.Session, orgID string, team *auth.Team) error {
	// FIXME: Orca doesn't seem to support all of the sync options that
	// eNZi offers. It only supports "group member selection" and not "user
	// search filter" for syncing team members. This means that it can't
	// support fancy things like nested group membership in Microsoft
	// Active Directory.
	memberSyncConfigForm := forms.MemberSyncOpts{
		EnableSync:         team.LdapDN != "", // Enable sync if LDAP DN is given.
		SelectGroupMembers: true,
		GroupDN:            team.LdapDN,
		GroupMemberAttr:    team.LdapMemberAttr,
	}

	if _, err := session.SetTeamMemberSyncConfig("id:"+orgID, "id:"+team.Id, memberSyncConfigForm); err != nil {
		return fmt.Errorf("unable to set team member sync config on auth provider: %s", err)
	}

	return nil
}

func (a *Authenticator) GetTeam(ctx *auth.Context, teamID string) (*auth.Team, error) {
	session := a.getSession(ctx.ClientCreds)

	teamResp, err := session.GetTeam("id:"+a.defaultOrgID, "id:"+teamID)
	if err != nil {
		apiErrs, ok := err.(*errors.APIErrors)
		if ok && apiErrs.HTTPStatusCode == http.StatusNotFound {
			return nil, auth.ErrTeamDoesNotExist
		}

		return nil, fmt.Errorf("unable to get team from auth provider: %s", err)
	}

	// Get team member usernames.
	// NOTE: We don't do this for the ListTeams handler because it would be
	// too expensive.
	membersResp, _, err := session.ListTeamMembers("id:"+a.defaultOrgID, "id:"+teamID, "", "", api.MaxPerPageLimit)
	if err != nil {
		return nil, fmt.Errorf("unable to list team members from auth provider: %s", err)
	}

	memberUsernames := make([]string, len(membersResp.Members))
	for i, memberResp := range membersResp.Members {
		memberUsernames[i] = memberResp.Member.Name
	}

	// Get team member sync config.
	memberSyncConfigResp, err := getTeamMemberSyncConfig(session, a.defaultOrgID, teamID)
	if err != nil {
		return nil, err
	}

	return &auth.Team{
		OrgId:       teamResp.OrgID,
		Name:        teamResp.Name,
		Id:          teamResp.ID,
		Description: teamResp.Description,
		// NOTE: We don't set the DiscoveredMembers field because eNZi
		// does does not differentiate between members that are added
		// manually vs those that are synced from LDAP. The eNZi API
		// server handles making sure that team membership cannot be
		// modified manually if it is supposed to be synced with LDAP.
		// The sync configuration is also enabled and controlled
		// separately.
		ManagedMembers: memberUsernames,
		// FIXME: Orca doesn't seem to support all of the sync options
		// that eNZi offers. It only supports "group member selection"
		// and not "user search filter" for syncing team members. This
		// means that it can't support fancy things like nested group
		// membership in Microsoft Active Directory.
		LdapDN:         memberSyncConfigResp.GroupDN,
		LdapMemberAttr: memberSyncConfigResp.GroupMemberAttr,
	}, nil
}

func getTeamMemberSyncConfig(session *client.Session, orgID, teamID string) (*responses.MemberSyncOpts, error) {
	// NOTE: this API call returns an error if the authenticated user does
	// not have 'admin' access to the team.. We can just ignore this error
	// if that's the case.
	memberSyncConfigResp, err := session.GetTeamMemberSyncConfig("id:"+orgID, "id:"+teamID)
	if err != nil {
		if apiErrs, ok := err.(*errors.APIErrors); ok && apiErrs.HTTPStatusCode == http.StatusForbidden {
			// The client user doesn't have access to view the sync
			// config for this team. It's okay to ignore this error
			// and return an empty value.
			return &responses.MemberSyncOpts{}, nil
		}

		return nil, fmt.Errorf("unable to get team member sync config from auth provider: %s", err)
	}

	return memberSyncConfigResp, nil
}

func (a *Authenticator) ListTeams(ctx *auth.Context) ([]*auth.Team, error) {
	session := a.getSession(ctx.ClientCreds)

	// The interface does not have pagination so ask for thedefault number.
	teamsResp, _, err := session.ListOrganizationTeams("id:"+a.defaultOrgID, "", api.MaxPerPageLimit)
	if err != nil {
		return nil, fmt.Errorf("unable to get teams from auth provider: %s", err)
	}

	teams := make([]*auth.Team, len(teamsResp.Teams))
	for i, teamResp := range teamsResp.Teams {
		// NOTE: We do not get LDAP member sync config for each team
		// because it would be too expensive.
		teams[i] = &auth.Team{
			OrgId:       teamResp.OrgID,
			Name:        teamResp.Name,
			Id:          teamResp.ID,
			Description: teamResp.Description,
		}
	}

	return teams, nil
}

func (a *Authenticator) DeleteTeam(ctx *auth.Context, team *auth.Team) error {
	session := a.getSession(ctx.ClientCreds)

	if err := session.DeleteTeam("id:"+a.defaultOrgID, "id:"+team.Id); err != nil {
		return fmt.Errorf("unable to delete team on auth provider: %s", err)
	}

	return nil
}
