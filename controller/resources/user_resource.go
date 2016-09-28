package resources

import (
	"github.com/docker/orca/auth"
)

// UserResourceRequest corresponds to the UserResource type ID
// UserResources are tied to a given username.
type UserResourceRequest struct {
	username string
}

func (r *UserResourceRequest) HasAccess(ctx *auth.Context) bool {
	if r.username == ctx.User.Username {
		return true
	}
	return false
}

func NewUserResource(username string) *UserResourceRequest {
	if username == "" {
		return nil
	}
	return &UserResourceRequest{
		username: username,
	}
}
