package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/lib/gc"
	"github.com/docker/dhe-deploy/shared/dtrutil"

	log "github.com/Sirupsen/logrus"
)

// main initiates a new GC within the job runner.
func main() {
	go setupSignals()

	// setup connects to rethink and etcd, performing basic sanity checking
	// before starting GC.
	setup, err := gc.NewSetup()
	if err != nil {
		log.WithField("error", err).Fatal("unable to start setup")
	}

	if len(os.Args) > 1 && os.Args[1] == "set-rw-mode" {
		if err = setup.SetReadWrite(); err != nil {
			log.WithField("error", err).Fatal("unable to set read write")
		}
		return
	}

	// connect to rethink prior to setting read-only just in case this
	// fails; downtime is minimized
	replicaID := os.Getenv(deploy.ReplicaIDEnvVar)
	db, err := dtrutil.GetRethinkSession(replicaID)
	if err != nil {
		log.WithField("error", err).Fatal("unable to get rethink session")
	}

	mode, err := setup.GetGCMode()
	if err != nil {
		log.WithField("error", err).Fatal("unable to get gc mode")
	}

	// set read only then begin GC
	if err := setup.SetReadOnly(); err != nil {
		log.WithField("error", err).Fatal("unable to set read only mode")
	}
	defer func() {
		if err := setup.SetReadWrite(); err != nil {
			log.WithField("error", err).Fatal("unable to set read write mode")
		}
	}()

	err = gc.NewGC(gc.Opts{
		Session: db,
		Config:  setup.GetConfig(),
		Mode:    mode,
	}).Run()
	if err != nil {
		log.WithField("error", err).Fatal("error running gc")
	}
}

func setupSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.WithField("signal", sig).Fatal("detected os signal; stopping GC")
}
