package proxy

import (
	"log"
	"sync"
	"time"

	"github.com/Microsoft/go-winio"
	datakit "github.com/docker/datakit/api/go-datakit"
	"github.com/docker/pinata/v1/pinataSockets"
	"golang.org/x/net/context"
)

// NewHTTPProxyRewriter returns a new HTTPProxyRewriter
func NewHTTPProxyRewriter() EnvRewriter {
	once := &sync.Once{}
	var client *datakit.Client

	return &HTTPProxyRewriter{
		clientSupplier: func(ctx context.Context) *datakit.Client {
			once.Do(func() {
				for {
					conn, err := winio.DialPipe(pinataSockets.GetDBSocketPath(), nil)
					if err != nil {
						log.Println("Failed to connect to the database", err)
						time.Sleep(100 * time.Millisecond)
						continue
					}

					c, err := datakit.NewClient(ctx, conn)
					if err != nil {
						log.Println("Failed to connect to the database", err)
						time.Sleep(100 * time.Millisecond)
						continue
					}

					client = c
					break
				}
			})

			return client
		},
	}
}
