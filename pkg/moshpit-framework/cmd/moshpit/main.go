package main

import (
	"os"

	"golang.org/x/net/context"

	"github.com/docker/dhe-deploy/pkg/moshpit-framework"
	"github.com/docker/dhe-deploy/pkg/moshpit-framework/builder"
	"github.com/docker/dhe-deploy/pkg/moshpit-framework/commands"
)

func Setup(ctx context.Context, setupConfig string) error {
	log := moshpit.LoggerFromCtx(ctx)
	log.Info("example setup function")
	return nil
}

type Job struct {
	ctx context.Context
}

func (j *Job) State() (commands.JobState, error) {
	log := moshpit.LoggerFromCtx(j.ctx)
	log.Info("example state function")
	return commands.JobState_SUCCESS, nil
}

func ClientRun(ctx context.Context, name, clientConfig string) (moshpit.Job, error) {
	log := moshpit.LoggerFromCtx(ctx)
	log.Info("example client run function")
	return &Job{ctx: ctx}, nil
}

func main() {
	app := builder.BuildMoshpit(Setup, ClientRun)
	app.Run(os.Args)
}
