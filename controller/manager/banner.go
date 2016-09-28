package manager

import (
	"github.com/docker/orca/auth"
)

const (
	BannerINFO = "INFO" // TODO Might not want low urgency messages to show up here...
	BannerWARN = "WARN"
	BannerCRIT = "CRIT"
)

type Banner struct {
	Level string `json:"level"`
	// Empty on OK level, otherwise the message to display in the banner
	Message string `json:"msg"`
}

func (m DefaultManager) GetBanner(ctx *auth.Context) []Banner {

	// TODO Refactor this to use channels/goroutines so we can parallelize and prioritize on the fly
	//      with the highest priority items at the front of the list

	banners := m.getLicenseBanners()
	banners = append(banners, m.getHABanners()...)
	if ctx.User != nil && ctx.User.Admin {
		banners = append(banners, m.getCABanners()...)
		banners = append(banners, m.getUpdateBanners()...)
	}

	// Show version skew banners to everyone (displayed during upgrade of an HA cluster)
	banners = append(banners, m.getVersionSkewBanners()...)

	// TODO Append other subsystems here...
	return banners
}
