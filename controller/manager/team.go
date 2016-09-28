package manager

import (
	"fmt"
	"github.com/docker/orca/auth"
)

var (
	ksTeams = datastoreVersion + "/teams"
)

func (m DefaultManager) Team(ctx *auth.Context, id string) (*auth.Team, error) {
	return m.GetAuthenticator().GetTeam(ctx, id)
}

func (m DefaultManager) Teams(ctx *auth.Context) ([]*auth.Team, error) {
	return m.GetAuthenticator().ListTeams(ctx)
}

func (m DefaultManager) SaveTeam(ctx *auth.Context, team *auth.Team) (string, error) {
	eventType, err := m.GetAuthenticator().SaveTeam(ctx, team)
	if err != nil {
		return "", err
	}
	m.logEvent(eventType, fmt.Sprintf("name=%s", team.Name), []string{"security"})
	return eventType, nil
}

func (m DefaultManager) DeleteTeam(ctx *auth.Context, team *auth.Team) error {
	if err := m.GetAuthenticator().DeleteTeam(ctx, team); err != nil {
		return err
	}
	m.logEvent("delete-team", fmt.Sprintf("name=%s", team.Name), []string{"security"})
	return nil
}

func (m DefaultManager) AddMemberToTeam(ctx *auth.Context, teamID, username string) error {
	if err := m.GetAuthenticator().AddTeamMember(ctx, teamID, username); err != nil {
		return err
	}
	m.logEvent("update-team-add-member", fmt.Sprintf("id=%s account=%s", teamID, username), []string{"security"})
	return nil
}

func (m DefaultManager) RemoveMemberFromTeam(ctx *auth.Context, teamID, username string) error {
	if err := m.GetAuthenticator().DeleteTeamMember(ctx, teamID, username); err != nil {
		return err
	}
	m.logEvent("update-team-remove-member", fmt.Sprintf("id=%s account=%s", teamID, username), []string{"security"})
	return nil
}
