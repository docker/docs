package flags

import "github.com/codegangsta/cli"

type ClientFlagValues struct {
	Server string
	Name   string
}

type FlagValues struct {
	Debug      bool
	ConfigFile string
	ConfigData string
	Client     ClientFlagValues
}

var (
	ConfigFileFlag = cli.StringFlag{
		Name:   "config-file",
		EnvVar: "MOSHPIT_CONFIG_FILE",
		Value:  "config.yml",
	}
	ConfigDataFlag = cli.StringFlag{
		Name:   "config-data",
		EnvVar: "MOSHPIT_CONFIG_DATA",
	}
	DebugFlag = cli.BoolFlag{
		Name:   "debug",
		EnvVar: "MOSHPIT_DEBUG",
	}
	ServerFlag = cli.StringFlag{
		Name:   "server",
		EnvVar: "MOSHPIT_SERVER",
	}
	NameFlag = cli.StringFlag{
		Name:   "name",
		EnvVar: "MOSHPIT_NAME",
	}
)

func WiredFlags() (map[string][]cli.Flag, *FlagValues) {
	values := FlagValues{}
	configFile := ConfigFileFlag
	configFile.Destination = &values.ConfigFile
	configData := ConfigDataFlag
	configData.Destination = &values.ConfigData
	debugFlag := DebugFlag
	debugFlag.Destination = &values.Debug
	serverFlag := ServerFlag
	serverFlag.Destination = &values.Client.Server
	nameFlag := NameFlag
	nameFlag.Destination = &values.Client.Name

	flags := map[string][]cli.Flag{
		"": {
			debugFlag,
		},
		"dropper": {
			configFile,
			configData,
		},
		"server": {
			configFile,
			configData,
		},
		"client": {
			serverFlag,
			nameFlag,
		},
	}
	return flags, &values
}
