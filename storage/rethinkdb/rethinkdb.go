package rethinkdb

import (
	"github.com/dancannon/gorethink"
	"github.com/docker/go-connections/tlsconfig"
)

var session *gorethink.Session

// Connection sets up a RethinkDB connection to the host (`host:port` format)
// using the CA .pem file provided at path `caFile` and the authKey
func Connection(caFile, host, authKey string) (*gorethink.Session, error) {
	tlsOpts := tlsconfig.Options{
		CAFile: caFile,
	}
	t, err := tlsconfig.Client(tlsOpts)
	if err != nil {
		return nil, err
	}
	return gorethink.Connect(
		gorethink.ConnectOpts{
			Address:   host,
			AuthKey:   authKey,
			TLSConfig: t,
		},
	)
}
