package tracking

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/controller/ctx"
	"github.com/docker/orca/controller/manager"
	"github.com/docker/orca/version"
	"github.com/timehop/go-mixpanel"
)

type Tracker struct {
	manager  manager.Manager
	excludes []string
}

func filterURI(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	return u.Path, nil
}

func NewTracker(m manager.Manager, excludes []string) *Tracker {
	return &Tracker{
		manager:  m,
		excludes: excludes,
	}
}

func (t *Tracker) LayerHandler(rc *ctx.OrcaRequestContext) (int, error) {
	skipTracking := false
	username := rc.Auth.User.Username
	r := rc.Request

	path, err := filterURI(r.RequestURI)
	if err != nil {
		log.Errorf("tracking path filter error: %s", err)
		return http.StatusInternalServerError, err
	}

	// check if excluded
	for _, e := range t.excludes {
		match, err := regexp.MatchString(e, path)
		if err != nil {
			log.Errorf("tracking exclude error: %s", err)
		}

		if match {
			skipTracking = true
			break
		}
	}

	if !t.manager.GetTrackingDisabled() && username != "" && path != "" && !skipTracking {
		tagParts := strings.Split(path, "/")
		tag := tagParts[1]
		// Hash the ID for anonymity
		id := fmt.Sprintf("%x", sha1.Sum([]byte(t.manager.ID())))

		evt := map[string]interface{}{
			"type":    "api",
			"path":    path,
			"version": version.FullVersion(),
			"tags":    strings.Join([]string{"api", tag, strings.ToLower(r.Method)}, ", "),
		}

		go func() {
			mp := &mixpanel.Mixpanel{
				Token:   manager.MixpanelToken,
				BaseUrl: manager.MixpanelUrl,
			}
			mp.Track(id, "event", evt)
		}()
	}
	return http.StatusOK, err
}
