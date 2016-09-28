package resources

import (
	"github.com/docker/orca/auth"
)

// TeamResourceRequest corresponds to the TeamResource type ID
// TeamResources are tied to a given team ID.
type TeamResourceRequest struct {
	teamID string
}

func (r *TeamResourceRequest) HasAccess(ctx *auth.Context) bool {
	for _, teamID := range append(ctx.User.ManagedTeams, ctx.User.DiscoveredTeams...) {
		if r.teamID == teamID {
			return true
		}
	}
	return false
}

func NewTeamResource(teamID string) *TeamResourceRequest {
	if teamID == "" {
		return nil
	}
	return &TeamResourceRequest{
		teamID: teamID,
	}
}
