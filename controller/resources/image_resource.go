package resources

import (
	"github.com/docker/orca/auth"
)

type ImageResourceRequest struct {
	ImageName string
}

// TODO: implement RBAC for Images
func (r *ImageResourceRequest) HasAccess(ctx *auth.Context) bool {
	return true
}

func NewImageResource(imageName string) *ImageResourceRequest {
	return &ImageResourceRequest{
		ImageName: imageName,
	}
}
