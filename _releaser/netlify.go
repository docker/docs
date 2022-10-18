package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	netlify "github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/plumbing/operations"
	"github.com/netlify/open-api/go/porcelain"
	netlifyctx "github.com/netlify/open-api/go/porcelain/context"
)

type NetlifyCmd struct {
	Remove NetlifyRemoveCmd `kong:"cmd,name=remove"`
	Deploy NetlifyDeployCmd `kong:"cmd,name=deploy"`
}

type netlifyGlobalFlags struct {
	SiteName string `kong:"name='site-name',env='NETLIFY_SITE_NAME'"`
}

type NetlifyRemoveCmd struct {
	netlifyGlobalFlags
}

func (s *NetlifyRemoveCmd) Run() error {
	siteName := cleanSiteName(s.SiteName)
	c := newNetlifyClient(getEnvOrSecret("NETLIFY_AUTH_TOKEN"))
	site, err := c.getSite(siteName)
	if err != nil {
		return fmt.Errorf("failed to get site %q: %w", siteName, err)
	}
	if site == nil {
		log.Printf("INFO: site %s already removed", siteName)
		return nil
	}
	return c.DeleteSite(c.ctx, site.ID)
}

type NetlifyDeployCmd struct {
	netlifyGlobalFlags
	PublishDir string `kong:"name='publish-dir',env='NETLIFY_PUBLISH_DIR'"`
}

func (s *NetlifyDeployCmd) Run() error {
	if ok, err := isDirEmpty(s.PublishDir); err != nil {
		return fmt.Errorf("cannot read publish dir %q: %w", s.PublishDir, err)
	} else if !ok {
		return fmt.Errorf("publish dir %q empty", s.PublishDir)
	}

	siteName := cleanSiteName(s.SiteName)
	c := newNetlifyClient(getEnvOrSecret("NETLIFY_AUTH_TOKEN"))

	site, err := c.CreateSite(c.ctx, &netlify.SiteSetup{
		Site: netlify.Site{
			AccountSlug: getEnvOrSecret("NETLIFY_ACCOUNT_SLUG"),
			Name:        siteName,
		},
	}, false)
	if err != nil {
		return fmt.Errorf("failed to create site %q: %w", siteName, err)
	}

	deploy, err := c.DeploySite(c.ctx, porcelain.DeployOptions{
		SiteID: site.ID,
		Dir:    s.PublishDir,
	})
	if err != nil {
		return fmt.Errorf("failed to deploy site %q (%s): %w", siteName, site.ID, err)
	}

	log.Printf("INFO: site %q (%s) deployed at %s\n", siteName, site.ID, deploy.URL)
	return nil
}

type netlifyClient struct {
	ctx netlifyctx.Context
	*porcelain.Netlify
}

func newNetlifyClient(authToken string) *netlifyClient {
	return &netlifyClient{
		ctx:     netlifyctx.WithAuthInfo(context.Background(), authInfo(authToken)),
		Netlify: porcelain.NewRetryableHTTPClient(strfmt.Default, 5),
	}
}

func (c *netlifyClient) getSite(siteName string) (*netlify.Site, error) {
	sites, err := c.ListSites(c.ctx, &operations.ListSitesParams{
		Context: c.ctx,
		Name:    &siteName,
	})
	if err != nil {
		return nil, err
	}
	for _, site := range sites {
		if site.Name == siteName {
			return site, nil
		}
	}
	return nil, nil
}

func authInfo(authToken string) runtime.ClientAuthInfoWriter {
	return runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		err := r.SetHeaderParam("User-Agent", "User-Agent: releaser/"+version)
		if err != nil {
			return fmt.Errorf("unable to set useragent header: %w", err)
		}
		err = r.SetHeaderParam("Authorization", "Bearer "+authToken)
		if err != nil {
			return fmt.Errorf("unable to set authorization header: %w", err)
		}
		return nil
	})
}

func cleanSiteName(siteName string) string {
	siteName = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(siteName, "-")
	siteName = regexp.MustCompile(`^-+\|-+$`).ReplaceAllString(siteName, "")
	return strings.ToLower(strings.TrimSuffix(siteName, "-"))
}

func isDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
