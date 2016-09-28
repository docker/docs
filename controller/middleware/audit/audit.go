package audit

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca"
	"github.com/docker/orca/controller/ctx"
	"github.com/docker/orca/controller/manager"
)

type Auditor struct {
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

func NewAuditor(m manager.Manager, excludes []string) *Auditor {
	return &Auditor{
		manager:  m,
		excludes: excludes,
	}
}

func (a *Auditor) LayerHandler(rc *ctx.OrcaRequestContext) (int, error) {
	// debug logging per request
	skipAudit := false

	// Typecast the context object and extract username
	username := rc.Auth.User.Username
	if username == "" {
		log.Errorf("unable to get username from context")
		return http.StatusInternalServerError, errors.New("unable to get username from context")
	}

	r := rc.Request
	path, err := filterURI(r.RequestURI)
	if err != nil {
		log.Errorf("audit path filter error: %s", err)
		return http.StatusInternalServerError, err
	}

	// check if excluded
	for _, e := range a.excludes {
		match, err := regexp.MatchString(e, path)
		if err != nil {
			log.Errorf("audit exclude error: %s", err)
			return http.StatusInternalServerError, err
		}

		if match {
			skipAudit = true
			break
		}
	}

	if username != "" && path != "" && !skipAudit {
		tagParts := strings.Split(path, "/")
		tag := tagParts[1]

		evt := &orca.Event{
			Type:       "api",
			Time:       time.Now(),
			Username:   username,
			Message:    path,
			RemoteAddr: r.RemoteAddr,
			Tags:       []string{"api", tag, strings.ToLower(r.Method)},
		}

		if err := a.manager.SaveEvent(evt); err != nil {
			log.Errorf("error saving event: %s", err)
			return http.StatusInternalServerError, err
		}
	} else {
		// Only send debug message if we're not logging the event (reduce redundant chatter)
		log.Debugf("%s?%s method=%s", r.URL.Path, r.URL.Query().Encode(), r.Method)
	}
	return http.StatusOK, nil
}
