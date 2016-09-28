package resources

import (
	"github.com/docker/orca/auth"
)

// AdminResourceRequest corresponds to the AdminResource type ID
// AdminResources perform a second paranoia check for admin status of the user
type AdminResourceRequest struct {
}

func (r *AdminResourceRequest) HasAccess(ctx *auth.Context) bool {
	// In a perfect world, this is unreachable
	if ctx.User.Admin {
		return true
	}
	return false
}

func NewAdminResource() *AdminResourceRequest {
	return &AdminResourceRequest{}
}
