package dtrutil

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy/shared/containers"

	"github.com/coreos/etcd/pkg/transport"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

// GetRethinkSession constructs a new connection to the rethink database
func GetRethinkSession(replicaID string) (*rethink.Session, error) {
	var session *rethink.Session

	err := Poll(2*time.Second, 15, func() error {
		tlsInfo := transport.TLSInfo{
			CertFile: containers.RethinkCertStore.CertPath(),
			KeyFile:  containers.RethinkCertStore.KeyPath(),
			CAFile:   containers.RethinkCACertStore.CertPath(),
		}

		t, err := transport.NewTransport(tlsInfo, time.Minute)
		if err != nil {
			return err
		}

		session, err = rethink.Connect(rethink.ConnectOpts{
			Addresses:     []string{containers.Rethinkdb.BridgeName(replicaID)},
			DiscoverHosts: false,
			MaxIdle:       5,
			MaxOpen:       10,
			TLSConfig:     t.TLSClientConfig,
		})

		if err != nil {
			logrus.Debugf("Failed to connect to rethink: %s", err)
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	logrus.Debug("Connected to rethink")
	return session, nil
}
