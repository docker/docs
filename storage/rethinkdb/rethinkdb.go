package rethinkdb

import (
	"time"

	"github.com/docker/go-connections/tlsconfig"
	"gopkg.in/dancannon/gorethink.v2"
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
// using the CA .pem file provided at path `caFile`
func Connection(tlsOpts tlsconfig.Options, host string) (*gorethink.Session, error) {
	t, err := tlsconfig.Client(tlsOpts)
	if err != nil {
		return nil, err
	}
	return gorethink.Connect(
		gorethink.ConnectOpts{
			Address:   host,
			TLSConfig: t,
		},
	)
}
