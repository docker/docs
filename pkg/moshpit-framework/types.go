package moshpit

import (
	"github.com/docker/dhe-deploy/pkg/moshpit-framework/commands"
	"golang.org/x/net/context"
)

type DropperConfig struct {
	Username          string
	Password          string
	HubUsername       string
	HubPassword       string
	HubRefreshToken   string
	UCPURL            string
	UCPCA             string
	UCPInsecureTLS    bool
	MoshpitImage      string
	ServerNode        string
	ServerHost        string
	NumClients        int
	ClientConstraints string
	Spread            bool
	NoCleanup         bool
	Debug             bool
}

type ServerConfig struct {
	ListenIP   string
	ListenPort int
	AlsoClient bool
	// these are passed onto the plugin to handle as strings
	Setup  interface{}
	Client interface{}
}

type Config struct {
	Server  ServerConfig
	Dropper DropperConfig
}

type SetupFunc func(context.Context, string) error
type ClientRunFunc func(context.Context, string, string) (Job, error)

type Job interface {
	State() (commands.JobState, error)
}
