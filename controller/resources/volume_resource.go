package resources

import (
	"github.com/docker/orca/auth"
)

type VolumeResourceRequest struct {
	VolumeName string
}

// TODO: implement RBAC for Volumes
func (r *VolumeResourceRequest) HasAccess(ctx *auth.Context) bool {
	return true
}

func NewVolumeResource(volumeName string) *VolumeResourceRequest {
	return &VolumeResourceRequest{
		VolumeName: volumeName,
	}
}
