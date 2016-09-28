package notary

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/client/transport"
	"github.com/docker/engine-api/types"
	kvstore "github.com/docker/libkv/store"
	"github.com/docker/notary/client"
	"github.com/docker/notary/trustpinning"
	"github.com/docker/notary/tuf/store"
	"github.com/docker/orca"
	"github.com/docker/orca/controller/ctx"
	"github.com/docker/orca/controller/manager"
	"github.com/docker/orca/controller/resources"
	"github.com/docker/orca/utils"
)

const (
	NotaryUrl              = "https://notary.docker.io"
	notaryDatastoreVersion = "notary" + "/v1"
	notaryRootFileName     = "root.json"
	notaryTSFileName       = "timestamp.json"
)

type tokenGetter func(username, accessType, namedStr string, registry orca.Registry) (string, error)

type GarantToken struct {
	Token string `json:"token"`
}

type NotaryMiddleware struct {
	manager  manager.Manager
	cacheDir string
	tg       tokenGetter
}

var ErrParsingChallenge = errors.New("Error parsing challenge response")

func NewNotaryMiddleware(manager manager.Manager, cacheDir string, tg tokenGetter) *NotaryMiddleware {
	return &NotaryMiddleware{
		manager:  manager,
		cacheDir: cacheDir,
		tg:       tg,
	}
}

func (mw NotaryMiddleware) getRegistry(namedStr string) (orca.Registry, error) {
	named, err := reference.ParseNamed(namedStr)
	if err != nil {
		return nil, err
	}
	hostname, _ := reference.SplitHostname(named)

	reg, err := mw.manager.Registry(hostname)
	if err != nil {
		return nil, err
	}
	return reg, nil
}

func (mw NotaryMiddleware) getRegTokenFromDTR(username, repoName string) (string, error) {
	// grab reg
	reg, err := mw.getRegistry(repoName)
	if err != nil {
		log.Debugf("NotaryMiddleware: Could not retrieve registry: %s", err)
		return "", err
	}
	// only need "pull" permission to view signed metadata for tags
	tok, err := mw.tg(username, "pull", repoName, reg)
	if err != nil {
		log.Debugf("NotaryMiddleware: Error calling tokenGetter: %s", err)
		return "", err
	}

	return tok, nil
}

func (mw NotaryMiddleware) getRegTransport(repoName string) (http.RoundTripper, error) {
	reg, err := mw.getRegistry(repoName)
	if err != nil {
		log.Debugf("NotaryMiddleware: Could not retrieve registry: %s", err)
		return nil, err
	}

	return reg.GetTransport(), nil
}

func (mw NotaryMiddleware) getNotaryTokenFromGarant(fullRepoName string, authConfig types.AuthConfig) (string, error) {
	// first get back a challenge to find the auth URL
	req, err := http.NewRequest("GET", NotaryUrl, nil)
	if err != nil {
		log.Debugf("NotaryMiddleware: Creating request failed with err: %s", err)
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Debugf("NotaryMiddleware: Getting challenge failed with err: %s", err)
		return "", err
	}

	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()

	// check the Www-Authenticate header to figure out where to request token from
	respHeader := resp.Header.Get("Www-Authenticate")
	// header should be of the form 'Bearer realm="{TOKEN_URL}",service="{SERVICE}"'
	parser := regexp.MustCompile(`Bearer realm="(.*)",service="(.*)"`)
	match := parser.FindStringSubmatch(respHeader)
	if len(match) < 3 {
		log.Debugf("NotaryMiddleware: FindStringSubmatch on %s failed with err: %s", respHeader, err)
		return "", ErrParsingChallenge
	}

	challengeUrl := match[1]
	service := match[2]

	uri := fmt.Sprintf("%s?scope=repository:%s:%s&service=%s", challengeUrl, fullRepoName, "pull", service)

	req, err = http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Debug("NotaryMiddleware: attempted to prepare invalid GET request for authentication")
		return "", err
	}

	// auth might be blank if no auth is required
	if authConfig.Username != "" {
		req.SetBasicAuth(authConfig.Username, authConfig.Password)
	} else if (authConfig != types.AuthConfig{}) {
		log.Debug("NotaryMiddleware: unhandled non-empty config")
		// TODO: implement support for other auth config fields
	}

	resp, err = client.Do(req)
	if err != nil {
		log.Debugf("NotaryMiddleware: Error fetching token: %s", uri)
		return "", err
	}

	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debugf("NotaryMiddleware: Error reading body: %s", uri)
		return "", err
	}

	var token GarantToken
	if err := json.Unmarshal(bts, &token); err != nil {
		log.Debugf("NotaryMiddleware: Error unmarshalling garant token: %s", uri)
		return "", err
	}

	return token.Token, nil
}

func (mw NotaryMiddleware) LayerHandler(rc *ctx.OrcaRequestContext) (int, error) {
	// get config flag to see whether to check notary
	if rc.RequiresNotary {
		var found bool
		var serviceReq *resources.ServiceResourceRequest
		serviceReq, found = rc.MainResource.(*resources.CRUDResourceRequest).LabelledResource.(*resources.ServiceResourceRequest)
		if !found {
			return http.StatusInternalServerError, errors.New(fmt.Sprintf("internal error: unable to locate service resource in context. this probably means the notary middleware is being called on an unsupported request type."))
		}

		named, err := reference.ParseNamed(utils.GetRepoFullName(serviceReq.ServiceSpec.TaskTemplate.ContainerSpec.Image))
		if err != nil {
			log.Debugf("NotaryMiddleware: Error parsing repo: %s", err)
			return http.StatusInternalServerError, err
		}

		// Parse the tag and image name before creating a Notary repository
		// Default to "latest" if the image name does not contain a tag
		tag := "latest"
		repoName := named.Name()

		if parsedName, ok := named.(reference.NamedTagged); ok {
			tag = parsedName.Tag()
			repoName = parsedName.Name()
		}
		log.Debugf("NotaryMiddleware: Parsed image name: %s, and using tag: %s", repoName, tag)

		hostname, _ := reference.SplitHostname(named)

		var notaryUrl, notaryToken string

		// create roundtripper with default CAs
		var rt http.RoundTripper
		rt = &http.Transport{}

		// need to detect if req is going to hub or dtr
		if !utils.IsHub(hostname) {
			if !*mw.manager.RequireContentTrustForDTR() {
				return http.StatusOK, nil
			}

			// notary and dtr use the same loadbalancer
			notaryUrl = "https://" + hostname
			regToken, err := mw.getRegTokenFromDTR(rc.Auth.User.Username, repoName)
			if err != nil {
				log.Debugf("NotaryMiddleware: Error getting token from DTR: %s", err)
				return http.StatusInternalServerError, err
			}

			// tok is a base64 encoded json string
			decodedBytes, err := base64.StdEncoding.DecodeString(regToken)
			if err != nil {
				log.Debugf("NotaryMiddleware: Error base64decoding token: %s", err)
				return http.StatusInternalServerError, err
			}

			var parsed map[string]string
			err = json.Unmarshal(decodedBytes, &parsed)
			if err != nil {
				log.Debugf("NotaryMiddleware: Error unmarshalling token: %s", err)
				return http.StatusInternalServerError, err
			}

			// need the registry auth for dtr
			notaryToken = parsed["registrytoken"]

			// use reg's transport since it has its own cert
			rt, err = mw.getRegTransport(repoName)
			if err != nil {
				log.Debugf("NotaryMiddleware: Error getting registry transport: %s", err)
				return http.StatusInternalServerError, err
			}
		} else {
			if !*mw.manager.RequireContentTrustForHub() {
				return http.StatusOK, nil
			}

			notaryUrl = NotaryUrl

			// X-Registry-Auth header is a base64 encoded json with username/password
			regAuthToken := rc.Request.Header.Get("x-registry-auth")
			authConfig := types.AuthConfig{}

			if regAuthToken != "" {
				// base64decode the header
				decodedBytes, err := base64.StdEncoding.DecodeString(regAuthToken)
				if err != nil {
					log.Debugf("NotaryMiddleware: Error base64decoding auth header: %s", err)
					return http.StatusInternalServerError, err
				}
				err = json.Unmarshal(decodedBytes, &authConfig)
				if err != nil {
					log.Debugf("NotaryMiddleware: Error json decoding auth header: %s", err)
					return http.StatusInternalServerError, err
				}
			}

			notaryToken, err = mw.getNotaryTokenFromGarant(repoName, authConfig)
			if err != nil {
				log.Debugf("NotaryMiddleware: Error getting token from Garant: %s", err)
				return http.StatusInternalServerError, err
			}
		}

		header := http.Header{}
		header.Set("Authorization", "Bearer "+notaryToken)
		modifier := transport.NewHeaderRequestModifier(header)
		rt = transport.NewTransport(rt, modifier)

		// try checking cache first for root (anchoring trust) and timestamp (preventing version rollback) of this repo
		cachedRootData, err := mw.manager.Datastore().Get(path.Join(notaryDatastoreVersion, repoName, notaryRootFileName))
		if err != nil && err != kvstore.ErrKeyNotFound {
			log.Debugf("NotaryMiddleware: Error retrieving key for root metadata from underlying datastore")
			return http.StatusInternalServerError, err
		}
		cachedTSData, err := mw.manager.Datastore().Get(path.Join(notaryDatastoreVersion, repoName, notaryTSFileName))
		if err != nil && err != kvstore.ErrKeyNotFound {
			log.Debugf("NotaryMiddleware: Error retrieving key for timestamp metadata from underlying datastore")
			return http.StatusInternalServerError, err
		}

		// Best effort write of root and timestamp files to local metadata cache: /tuf/expanded/repo/name/.../metadata
		localMetaCacheDir := filepath.Join(mw.cacheDir, "tuf", filepath.FromSlash(repoName), "metadata")
		if err := os.MkdirAll(localMetaCacheDir, 0700); err != nil {
			log.Debugf("NotaryMiddleware: Error creating local cache directory")
		}
		if cachedRootData != nil {
			if err := ioutil.WriteFile(filepath.Join(localMetaCacheDir, notaryRootFileName), cachedRootData.Value, 0644); err != nil {
				log.Debugf("NotaryMiddleware: Error writing cached root metadata to local cache")
			}
		}

		if cachedTSData != nil {
			if err := ioutil.WriteFile(filepath.Join(localMetaCacheDir, notaryTSFileName), cachedTSData.Value, 0644); err != nil {
				log.Debugf("NotaryMiddleware: Error writing cached timestamp metadata to local cache")
			}
		}

		c, err := client.NewNotaryRepository(
			mw.cacheDir,
			repoName,
			notaryUrl,
			rt,
			nil,
			trustpinning.TrustPinConfig{},
		)
		if err != nil {
			log.Debugf("NotaryMiddleware: Creating notary repo returned error: %s", err)
			return http.StatusInternalServerError, err
		}

		trustData, err := c.GetAllTargetMetadataByName(tag)
		if err != nil {
			log.Debugf("NotaryMiddleware: Err getting targets: %s", err)
			// If notary is offline, we can surface without leakage
			// Else we hide if the repo or tag does not exist or if the user is unauthorized
			if _, ok := err.(store.ErrOffline); ok {
				return http.StatusInternalServerError, err
			}
			return http.StatusNotFound, fmt.Errorf("image or trust data does not exist for %s:%s", repoName, tag)
		}

		// If we successfully connected to notary and (possibly) updated, update the root and timestamp files for future use
		rootFileBytes, err := ioutil.ReadFile(filepath.Join(localMetaCacheDir, notaryRootFileName))
		// If the error is not nil, we somehow connected to notary but had our local cache corrupted; error out
		if err != nil {
			return http.StatusInternalServerError, err
		}
		if err := mw.manager.Datastore().Put(path.Join(notaryDatastoreVersion, repoName, notaryRootFileName), rootFileBytes, nil); err != nil {
			return http.StatusInternalServerError, err
		}
		tsFileBytes, err := ioutil.ReadFile(filepath.Join(localMetaCacheDir, notaryTSFileName))
		if err != nil {
			return http.StatusInternalServerError, err
		}
		if err := mw.manager.Datastore().Put(path.Join(notaryDatastoreVersion, repoName, notaryTSFileName), tsFileBytes, nil); err != nil {
			return http.StatusInternalServerError, err
		}

		trustDataMap := make(map[string]client.Target)
		for _, targetWithRole := range trustData {
			trustDataMap[targetWithRole.Role] = targetWithRole.Target
		}
		// First check if targets/releases has signed metadata for this role, then targets
		// If neither of those roles has signed, then error
		var digestBytes []byte
		if releasesData, ok := trustDataMap["targets/releases"]; ok {
			digestBytes = releasesData.Hashes["sha256"]
		} else if targetsData, ok := trustDataMap["targets"]; ok {
			digestBytes = targetsData.Hashes["sha256"]
		} else {
			return http.StatusNotFound, fmt.Errorf("image or trust data does not exist for %s:%s", repoName, tag)
		}
		serviceReq.ServiceSpec.TaskTemplate.ContainerSpec.Image = fmt.Sprintf("%s@sha256:%s", serviceReq.ServiceSpec.TaskTemplate.ContainerSpec.Image, hex.EncodeToString(digestBytes))
	}
	return http.StatusOK, nil
}
