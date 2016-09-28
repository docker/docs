package main

import (
	"net"
	"os"
	"os/exec"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy/shared/containers"
)

var interval = 15 * time.Second

func main() {
	log.Info("dnsfix daemon starting")
	namesToWatch := []string{}
	for _, container := range []containers.DTRContainer{
		containers.APIServer,
		containers.Registry,
		containers.NotaryServer,
		// these were also added to nginx's health check, so now we have to poll them here too
		containers.Rethinkdb,
		containers.Etcd,
	} {
		namesToWatch = append(namesToWatch, container.BridgeNameLocalReplica())
	}

	ipCache := map[string]net.IPAddr{}

	for range time.Tick(interval) {
		for _, name := range namesToWatch {
			// avoid doing too many requests
			addr, err := net.ResolveIPAddr("ip4", name)
			if err != nil {
				log.WithFields(log.Fields{
					"name":  name,
					"error": err,
				}).Warn("DNS lookup failed.")
				continue
			}
			if old, ok := ipCache[name]; !ok {
				ipCache[name] = *addr
			} else if !old.IP.Equal(addr.IP) {
				log.WithFields(log.Fields{
					"name": name,
					"old":  old,
					"new":  addr,
				}).Info("DNS name change detected. Reloading nginx.")

				reloadCmd := exec.Command("nginx", "-s", "reload")
				reloadCmd.Stdout = os.Stdout
				reloadCmd.Stderr = os.Stderr
				if err := reloadCmd.Run(); err != nil {
					log.WithField("error", err).Warn("error running nginx reload")
				}
			}
			ipCache[name] = *addr
		}
	}
}
