package manager

import (
	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/orca"
	"github.com/docker/orca/auth"
)

func (m DefaultManager) ListImages(all bool, f filters.Args) ([]types.Image, error) {
	return m.client.ImageList(context.TODO(), types.ImageListOptions{
		All:     all,
		Filters: f,
	})
}

func (m DefaultManager) ListUserImages(ctx *auth.Context, all bool, f filters.Args) ([]types.Image, error) {
	allImages, err := m.ListImages(all, f)
	if err != nil {
		return nil, err
	}

	acct := ctx.User

	// return immediately for admins
	if acct.Admin {
		return allImages, nil
	}

	access, err := m.GetAccess(ctx)
	if err != nil {
		return nil, err
	}

	userImages := []types.Image{}

	for _, i := range allImages {
		hasLabel := false

		// check labels for access
		for k, v := range i.Labels {
			if k != orca.UCPAccessLabel {
				continue
			}

			hasLabel = true
			if lvl, ok := access[v]; ok && lvl > auth.None {
				userImages = append(userImages, i)
			}
		}

		// by default if an image does not have a label it's public
		if !hasLabel {
			userImages = append(userImages, i)
		}
	}

	return userImages, nil
}
