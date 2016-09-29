package proxy

import (
	"log"
	"time"

	datakit "github.com/docker/datakit/api/go-datakit"
	"github.com/docker/pinata/v1/pinataSockets"
	"golang.org/x/net/context"
)

// NewHTTPProxyRewriter returns a new HTTPProxyRewriter
func NewHTTPProxyRewriter() *HTTPProxyRewriter {
	ctx := context.Background()

	for {
		c, err := datakit.Dial(ctx, "unix", pinataSockets.GetDBSocketPath())
		if err != nil {
			log.Println("Failed to connect to database", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		return &HTTPProxyRewriter{
			clientSupplier: func(ctx context.Context) *datakit.Client {
				return c
			},
		}
	}

}
