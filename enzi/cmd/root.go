package cmd

import (
	"crypto/tls"
	"fmt"

	log "github.com/Sirupsen/logrus"
	rethink "gopkg.in/dancannon/gorethink.v2"
	// To register job commands.
	"github.com/codegangsta/cli"
	"github.com/docker/orca/enzi/util"
)

// Root is the root command for eNZi.
var Root = cli.NewApp()

var (
	dbAddrs       = []string{"db.enzi"}
	debug         bool
	jsonLogFormat bool
)

func init() {
	Root.Name = "enzi"
	Root.HelpName = "enzi"
	Root.Usage = "An Auth Service for Docker Datacenter"
	Root.Version = "1.0.1"
	Root.Before = func(ctx *cli.Context) error {
		if debug {
			log.SetLevel(log.DebugLevel)
		}
		if jsonLogFormat {
			log.SetFormatter(&log.JSONFormatter{})
		}

		return nil
	}

	dbAddrsWrapper := &stringSlice{
		vals: &dbAddrs,
	}

	Root.Flags = []cli.Flag{
		cli.GenericFlag{
			Name:  "db-addr",
			Value: dbAddrsWrapper,
			Usage: "address (host[:port]) of database (may be specified multiple times)",
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "enable debug logs",
			Destination: &debug,
		},
		cli.BoolFlag{
			Name:        "jsonlog",
			Usage:       "format logs as JSON",
			Destination: &jsonLogFormat,
		},
	}
}

// stringSlice implements cli.Generic
// vals can be set to a default slice of strings but it is replaced with a
// new string slice as soon as a value is set. Any additional values are
// appended to the slice.
type stringSlice struct {
	vals    *[]string
	changed bool
}

func (s *stringSlice) String() string {
	return fmt.Sprintf("%v", *(s.vals))
}

func (s *stringSlice) Set(val string) error {
	if s.changed {
		*s.vals = append(*s.vals, val)
	} else {
		*s.vals = []string{val}
		s.changed = true
	}

	return nil
}

// GetTLSConfig loads a TLS config from the default location and uses the
// given clientAuth policy.
func GetTLSConfig(clientAuth tls.ClientAuthType) *tls.Config {
	tlsConfig, err := util.GetTLSConfig(clientAuth)
	if err != nil {
		log.Fatalf("unable to get TLS config: %s", err)
	}

	return tlsConfig
}

// GetDBSession connects to the RethinkDB cluster.
func GetDBSession(tlsConfig *tls.Config) *rethink.Session {
	log.Debugf("connecting to DB Addrs: %v", dbAddrs)

	dbSession, err := util.GetDBSession(dbAddrs, tlsConfig)
	if err != nil {
		log.Fatalf("unable to connect to database cluster: %s", err)
	}

	return dbSession
}
