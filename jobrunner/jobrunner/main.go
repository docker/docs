package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/WatchBeam/clock"
	"github.com/codegangsta/cli"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/jobrunner/jobrunner/jobrunner_config"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/worker"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	"github.com/docker/distribution/context"
)

type FlagValues struct {
	Debug bool
}

var (
	DebugFlag = cli.BoolFlag{
		Name:   "debug",
		EnvVar: "JOBRUNNER_DEBUG",
	}
)

func WiredFlags() (map[string][]cli.Flag, *FlagValues) {
	values := FlagValues{}
	debugFlag := DebugFlag
	debugFlag.Destination = &values.Debug

	flags := map[string][]cli.Flag{
		"": {
			debugFlag,
		},
		"worker": {},
	}
	return flags, &values
}

func main() {
	app := cli.NewApp()
	app.Name = "jobrunner"
	app.Usage = "Run periodic jobs in a symmetric HA system"
	flagMap, flagValues := WiredFlags()
	logrus.SetFormatter(new(logrus.JSONFormatter))
	level := logrus.InfoLevel
	app.Flags = flagMap[""]
	before := func(c *cli.Context) error {
		if flagValues.Debug {
			level = logrus.DebugLevel
		}
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:   "worker",
			Usage:  "Run a jobrunner worker",
			Flags:  flagMap["dropper"],
			Before: before,
			Action: func(c *cli.Context) {
				ctx := context.Background()
				// TODO: set log level in context
				//ctx := moshpit.MakeCtx("worker", level)
				log := context.GetLogger(ctx)
				var err error

				replicaID := os.Getenv(deploy.ReplicaIDEnvVar)
				dbSession, err := dtrutil.GetRethinkSession(replicaID)
				if err != nil {
					log.Fatal(err)
				}

				schemaMgr := schema.NewJobrunnerManager(deploy.JobrunnerDBName, dbSession)

				w := worker.New(ctx, schemaMgr, os.Getenv(deploy.ReplicaIDEnvVar), jobrunner_config.RegisteredActions, &clock.DefaultClock{}, nil, nil)

				log.Info("starting worker")
				err = w.Run()
				if err != nil {
					log.Fatal(err)
				}
			},
		},
	}
	app.Run(os.Args)
}
