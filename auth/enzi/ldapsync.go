package enzi

import (
	"github.com/docker/orca/auth"
)

// Sync is a no-op since syncing is handleded directly through the provider.
func (a *Authenticator) Sync(_ *auth.Context, _, _ bool) error {
	return nil
}

// LastSyncStatus is a no-op since sync job status is retreived directly from
// the provider now.
func (a *Authenticator) LastSyncStatus(_ *auth.Context) string {
	return ""
}
