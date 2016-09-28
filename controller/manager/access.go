package manager

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/auth"
)

// GetAccess builds a map of label to access using the highest access
// for each team that the account is a member
func (m DefaultManager) GetAccess(ctx *auth.Context) (map[string]auth.Role, error) {
	access := map[string]auth.Role{}

	log.Debugf("team membership: user=%s managed teams=%d discovered teams=%d", ctx.User.Username, len(ctx.User.ManagedTeams), len(ctx.User.DiscoveredTeams))

	teams := append(ctx.User.ManagedTeams, ctx.User.DiscoveredTeams...)
	for _, teamID := range teams {
		team, err := m.GetAuthenticator().GetTeam(ctx, teamID)
		if err != nil {
			log.Info("Stale team membership: %s - %s", teamID, err)
			continue
		}
		log.Debugf("checking access for team: name=%s", team.Name)
		lists, err := m.AccessListsForTeam(team.Id)
		if err != nil {
			return nil, err
		}

		// access role for label
		for _, list := range lists {
			role := list.Role
			if v, ok := access[list.Label]; ok && list.Role < v {
				continue
			}

			access[list.Label] = role
		}
	}

	return access, nil
}
