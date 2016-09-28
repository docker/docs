package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/distribution/reference"
	"github.com/gorilla/mux"

	"github.com/docker/orca"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/controller/ctx"
	"github.com/docker/orca/controller/manager"
	"github.com/docker/orca/utils"
)

type containerCreateData struct {
	Image string `json:"Image"`
}

// TODO(alexmavr): Rewrite most of these methods using the new pipeline

func (a *Api) swarmRegistryImagesPush(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	a.swarmRegistryRedirect(w, rc, rc.PathVars["name"], "push,pull")
}

func (a *Api) swarmRegistryImagesCreate(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	repo, err := utils.GetRepoFromFormVars(rc.Request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.swarmRegistryRedirect(w, rc, repo, "pull")
}

func (a *Api) swarmRegistryRedirect(w http.ResponseWriter, rc *ctx.OrcaRequestContext, repo string, accessType string) {
	a.registryRedirect(a.swarmClassicURL, w, rc, repo, accessType)
}

func (a *Api) engineRegistryRedirect(w http.ResponseWriter, rc *ctx.OrcaRequestContext, repo string, accessType string) {
	a.registryRedirect(a.engineProxyURL, w, rc, repo, accessType)
}

func (a *Api) registryRedirect(redirectURL string, w http.ResponseWriter, rc *ctx.OrcaRequestContext, repo string, accessType string) {
	var err error
	req := rc.Request

	survivable := func(err error) bool {
		switch err {
		case manager.ErrRegistryDoesNotExist:
			log.Debugf("Couldn't find registry")
			return true
		case auth.ErrUnauthorized:
			log.Warnf("Not authorized to contact registry for this user")
			return true
		case auth.ErrUnknown:
			log.Warnf("Unknown error occurred while contacting registry")
			return true
		}
		log.Errorf("There was an error while attempting to get a token: %s", err)
		return false
	}

	ok, err := a.replaceRegistryHeader(req, repo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ok {
		// Try to get a garant token so that we can push to the registry.
		// If we're not able to get a token (e.g. registry isn't registered to
		// to UCP, couldn't contact the registry, the registry says this user
		// is unauthorized, etc.), just send the request through in case the
		// user has access by some other means.
		header, err := a.getAuthHeaderForRepo(rc.Auth.User.Username, accessType, repo)

		if err != nil && !survivable(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if header != "" {
			req.Header.Set("X-Registry-Auth", header)
		}
	}
	req.URL, err = url.ParseRequestURI(redirectURL)
	if err != nil {
		log.Debugf("problem with parsing request.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.fwd.ServeHTTP(w, req)
}

func (a *Api) getRepoFromRequest(req *http.Request) (string, error) {
	vars := mux.Vars(req)
	return vars["name"], nil
}

func (a *Api) replaceRegistryHeader(req *http.Request, repo string) (bool, error) {
	data := req.Header.Get("X-Registry-Auth")

	// no need to get a token if this is for a non-registry operation (e.g.
	// "docker load")
	if repo == "" {
		return false, nil
	}

	// always OK to replace the auth header in the unlikely case that
	// it doesn't exist.
	if data == "" {
		return true, nil
	}

	// docker client uses URL encoding and not standard encoding
	authJSON := base64.NewDecoder(base64.URLEncoding, strings.NewReader(data))
	var authCreds *orca.AuthConfig
	if err := json.NewDecoder(authJSON).Decode(&authCreds); err != nil {
		return false, err
	}

	// Split the repo apart to try and get the hostname, if it exists (which it probably
	// won't for pushes/pulls to the Hub).  Unfortunately we'll interpret
	// the first part of the namespace as the hostname.  In most cases, this is fine
	// (we'll fail to find the registry and won't bother trying to get a token), however
	// if the namespace on the Hub is named the same as the local registry, we'll
	// substitute in a token for the local registry and access to the Hub will fail.
	//
	// This is an extremely unlikely corner case.
	named, err := reference.ParseNamed(repo)
	if err != nil {
		return false, err
	}
	hostname, _ := reference.SplitHostname(named)

	// don't replace the auth header if there are already existing
	// auth credentials for this registry.
	if authCreds.Hostname == hostname {
		return false, nil
	}
	return true, nil
}

func (a *Api) getAuthHeaderForRepo(username, accessType, namedStr string) (string, error) {
	named, err := reference.ParseNamed(namedStr)
	if err != nil {
		return "", err
	}
	hostname, reponame := reference.SplitHostname(named)

	reg, err := a.manager.Registry(hostname)
	if err != nil {
		return "", err
	}

	// Try to get a token
	if reg != nil {
		token, err := reg.GetAuthToken(username, accessType, hostname, reponame)
		if err != nil {
			if err == auth.ErrUnauthorized {
				log.Warnf("User '%s' was not authorized to access '%s'", username, namedStr)
			}
			return "", err
		}
		tokenstr := fmt.Sprintf("{\"registrytoken\": \"%s\"}", token)
		return base64.StdEncoding.EncodeToString([]byte(tokenstr)), nil
	}
	return "", nil
}

func (a *Api) getAuthHeaderForGUN(username, accessType, namedStr string, reg orca.Registry) (string, error) {
	named, err := reference.ParseNamed(namedStr)
	if err != nil {
		return "", err
	}
	hostname, _ := reference.SplitHostname(named)

	token, err := reg.GetAuthToken(username, accessType, hostname, named.Name())
	if err != nil {
		if err == auth.ErrUnauthorized {
			log.Warnf("User '%s' was not authorized to access '%s'", username, namedStr)
		}
		return "", err
	}
	tokenstr := fmt.Sprintf("{\"registrytoken\": \"%s\"}", token)
	return base64.StdEncoding.EncodeToString([]byte(tokenstr)), nil
}
