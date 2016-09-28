package resources

import (
	log "github.com/Sirupsen/logrus"

	"github.com/docker/orca/auth"
	"github.com/docker/orca/controller/manager"
)

type ResourceType int

const (
	UnknownResource ResourceType = iota
	AdminResource
	UserResource
	TeamResource
	ContainerResource
	ImageResource
	NetworkResource
	VolumeResource
	ApplicationResource
	ConsoleSessionResource
	ContainerLogsResource
	PublicResource
)

type ResourceAction int

const (
	ActionView ResourceAction = iota
	ActionRestricted
	ActionFull
	ActionCreate // Create is a special action as the minimum role is defined
	// By the labels passed in the underlying request struct
)

// The ResourceRequest interface is explicitly casted by the handlers
type ResourceRequest interface {
	HasAccess(*auth.Context) bool
}

// A CRUDResourceRequest represents a Resource whose access control
// is managed exclusively by its Owner and Team labels
type CRUDResourceRequest struct {
	LabelledResource
	Action        ResourceAction
	effectiveRole auth.Role
	manager       manager.Manager
}

type LabelledResource interface {
	TeamLabel() (string, error)
	OwnerLabel() (string, error)
}

func (r *CRUDResourceRequest) GetEffectiveRole() auth.Role {
	return r.effectiveRole
}

func (r *CRUDResourceRequest) HasAccess(ctx *auth.Context) bool {
	switch r.Action {
	case ActionCreate:
		return r.hasAccessFromConfig(ctx, auth.RestrictedControl)
	case ActionRestricted:
		return r.hasAccessFromConfig(ctx, auth.RestrictedControl)
	case ActionView:
		return r.hasAccessFromConfig(ctx, auth.View)
	case ActionFull:
		return r.hasAccessFromConfig(ctx, auth.FullControl)
	default:
		return false
	}
}

func (r *CRUDResourceRequest) hasAccessFromConfig(ctx *auth.Context, minRole auth.Role) bool {
	r.effectiveRole = ctx.User.Role

	// Get the team label
	teamLabel, err := r.LabelledResource.TeamLabel()
	if err != nil {
		log.Error(err)
		return false
	}

	if teamLabel != "" {
		// Compare the team label of the request with the user's access labels
		access, err := r.manager.GetAccess(ctx)
		if err != nil {
			log.Error(err)
			return false
		}
		role, ok := access[teamLabel]
		if !ok {
			// Block access to a user with an unknown access label
			return false
		}
		log.Infof("Found team label match, role: %d", role)
		r.effectiveRole = role
		return role >= minRole
	}

	// Create operations don't need to check the owner label of the underlying API request
	// Immediately check if the user's default role is enough to satisfy the minimum
	if r.Action == ActionCreate {
		return ctx.User.Role >= minRole
	}

	// No team label match, check the owner label
	ownerLabel, err := r.LabelledResource.OwnerLabel()
	if err != nil {
		log.Error(err)
		return false
	}

	if ownerLabel != ctx.User.Username {
		return false
	}
	if ctx.User.Role >= minRole {
		return true
	}
	return false
}
