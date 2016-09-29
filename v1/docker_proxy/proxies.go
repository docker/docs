package proxy

import (
	"fmt"
	"log"
	"regexp"

	datakit "github.com/docker/datakit/api/go-datakit"
	"github.com/docker/engine-api/types/container"
	"golang.org/x/net/context"
)

const dir = "com.docker.driver.amd64-linux"

// EnvRewriter rewrites environment variables
type EnvRewriter interface {
	RewriteEnv(config *container.Config)
}

// HTTPProxyRewriter uses the database to read proxy variables
type HTTPProxyRewriter struct {
	clientSupplier func(ctx context.Context) *datakit.Client
}

// RewriteEnv rewrites the environment in a given configuration
func (p *HTTPProxyRewriter) RewriteEnv(config *container.Config) {
	if config == nil {
		return
	}

	if len(config.Env) > 0 {
		r := regexp.MustCompile("(?i)^http_proxy.*|^https_proxy.*|^no_proxy.*")
		// Check if proxy is already set
		for _, s := range config.Env {
			if r.MatchString(s) {
				// user has already set proxies, don't override
				return
			}
		}
	}

	ctx := context.Background()

	client := p.clientSupplier(ctx)

	sha, err := datakit.Head(ctx, client, "master")
	if err != nil {
		log.Printf("Failed to discover HEAD of master %#v", err)
		return
	}

	snap := datakit.NewSnapshot(ctx, client, datakit.COMMIT, sha)

	httpProxy, err := snap.Read(ctx, []string{dir, "proxy", "http"})
	if err != nil {
		log.Printf("Failed to read proxies/http from snaphshot %#v", err)
		return
	}

	if httpProxy != "" {
		config.Env = append(config.Env,
			fmt.Sprintf("HTTP_PROXY=%s", httpProxy),
			fmt.Sprintf("http_proxy=%s", httpProxy),
		)
	}

	httpsProxy, err := snap.Read(ctx, []string{dir, "proxy", "https"})
	if err != nil {
		log.Printf("Failed to read proxies/https from snaphshot %#v", err)
		return
	}

	if httpsProxy != "" {
		config.Env = append(config.Env,
			fmt.Sprintf("HTTPS_PROXY=%s", httpsProxy),
			fmt.Sprintf("https_proxy=%s", httpsProxy),
		)
	}

	noProxy, err := snap.Read(ctx, []string{dir, "proxy", "exclude"})
	if err != nil {
		log.Printf("Failed to read proxies/exclude from snaphshot %#v", err)
		return
	}

	if noProxy != "" {
		config.Env = append(config.Env,
			fmt.Sprintf("no_proxy=%s", noProxy),
		)
	}
}
