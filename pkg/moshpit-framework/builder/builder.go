package builder

import (
	"fmt"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/dhe-deploy/pkg/moshpit-framework"
	"github.com/docker/dhe-deploy/pkg/moshpit-framework/dropper"
	"github.com/docker/dhe-deploy/pkg/moshpit-framework/flags"
	"github.com/docker/dhe-deploy/pkg/moshpit-framework/harness"
	"gopkg.in/yaml.v2"
)

// This can be expanded to take in a custom protobufs struct for the MoshConfig,
// some custom server code and some custom client code
// then this can become a general purpose load testing framework
func BuildMoshpit(setupFunc moshpit.SetupFunc, clientRunFunc moshpit.ClientRunFunc) *cli.App {
	app := cli.NewApp()
	app.Name = "moshpit"
	app.Usage = "Load-test network things"
	flagMap, flagValues := flags.WiredFlags()
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
			Name:   "dropper",
			Flags:  flagMap["dropper"],
			Before: before,
			Action: func(c *cli.Context) {
				ctx := moshpit.MakeCtx("dropper", level)
				log := moshpit.LoggerFromCtx(ctx)
				var err error
				configBytes := []byte(flagValues.ConfigData)
				if len(configBytes) == 0 {
					configBytes, err = ioutil.ReadFile(flagValues.ConfigFile)
					if err != nil {
						log.Fatal(err.Error())
					}
				}
				config := moshpit.Config{}
				err = yaml.Unmarshal(configBytes, &config)
				if err != nil {
					log.Fatal(err.Error())
				}
				d := dropper.New(ctx, config)
				err = d.Execute()
				if err != nil {
					log.Fatal(err.Error())
				}
			},
		},
		{
			Name:   "server",
			Flags:  flagMap["server"],
			Before: before,
			Action: func(c *cli.Context) {
				ctx := moshpit.MakeCtx("server", level)
				log := moshpit.LoggerFromCtx(ctx)
				var err error
				configBytes := []byte(flagValues.ConfigData)
				if len(configBytes) == 0 {
					configBytes, err = ioutil.ReadFile(flagValues.ConfigFile)
					if err != nil {
						log.Fatal(err.Error())
					}
				}
				config := moshpit.Config{}
				err = yaml.Unmarshal(configBytes, &config)
				if err != nil {
					log.Fatal(err.Error())
				}

				// TODO: implement tls between client and server
				if !config.Server.AlsoClient {
					err = harness.Server(ctx, config.Server, setupFunc, clientRunFunc)
					if err != nil {
						log.Fatal(err.Error())
					}
				} else {
					go func() {
						err = harness.Server(ctx, config.Server, setupFunc, clientRunFunc)
						if err != nil {
							log.WithField("error", err).Warn("Server existed with error")
						}
					}()
					ctx2 := moshpit.MakeCtx("client", level)
					err := harness.Client(ctx2, fmt.Sprintf("127.0.0.1:%d", config.Server.ListenPort), "client", clientRunFunc)
					if err != nil {
						log.Fatal(err.Error())
					}
				}

			},
		},
		{
			Name:   "client",
			Flags:  flagMap["client"],
			Before: before,
			Action: func(c *cli.Context) {
				ctx := moshpit.MakeCtx(flagValues.Client.Name, level)
				log := moshpit.LoggerFromCtx(ctx)
				// TODO: implement tls between client and server
				err := harness.Client(ctx, flagValues.Client.Server, flagValues.Client.Name, clientRunFunc)
				if err != nil {
					log.Fatal(err.Error())
				}
			},
		},
	}
	return app
}
