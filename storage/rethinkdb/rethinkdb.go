package rethinkdb

import (
	"time"

	"github.com/dancannon/gorethink"
	"github.com/docker/go-connections/tlsconfig"
)

var session *gorethink.Session

// Timing can be embedded into other gorethink models to
// add time tracking fields
type Timing struct {
	CreatedAt time.Time `gorethink:"created_at"`
	UpdatedAt time.Time `gorethink:"updated_at"`
	DeletedAt time.Time `gorethink:"deleted_at"`
}

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
