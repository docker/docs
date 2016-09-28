package resources

import (
	"github.com/docker/orca/auth"
)

// PublicResourceRequest corresponds to PublicResource
// PublicResources are accessible by every authenticated user of the system
type PublicResourceRequest struct {
}

func (r *PublicResourceRequest) HasAccess(ctx *auth.Context) bool {
	return true
}

func NewPublicResource() *PublicResourceRequest {
	return &PublicResourceRequest{}
}
