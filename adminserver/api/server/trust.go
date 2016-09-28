package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
	"github.com/docker/dhe-deploy/bootstrap"
	"github.com/docker/dhe-deploy/garant/authn"
	"github.com/docker/dhe-deploy/manager/schema"
	registryclient "github.com/docker/dhe-deploy/registry/client"
	"github.com/docker/dhe-deploy/shared/containers"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/registry/client/transport"
	"github.com/docker/notary"
	notaryclient "github.com/docker/notary/client"
	"github.com/docker/notary/trustpinning"
	enziresponses "github.com/docker/orca/enzi/api/responses"
	"github.com/emicklei/go-restful"
)

type roundTripBuilder func(string) (http.RoundTripper, error)

const notaryDataPath = "/var/local/dtr/notary-client"

func (a *APIServer) createNotaryTransportBuilder(ctx context.Context, user *authn.User, namespace *enziresponses.Account, repo *schema.Repository, grantedAccessLevel string) (roundTripBuilder, error) {
	// add the notary cert and the regular CA cert, so we can speak to notary server
	// and also the garant server
	trans, err := bootstrap.GetNotaryTransport()
	if err != nil {
		return nil, err
	}

	opts := registryclient.Opts{
		a.settingsStore,
		a.metadataMgr,
		a.repoMgr,
	}
	registryClient, err := registryclient.NewClient(ctx, opts)
	if err != nil {
		return nil, err
	}

	return func(repoFullName string) (http.RoundTripper, error) {
		jwt, err := registryClient.CreateJWT(user, repoFullName, grantedAccessLevel)
		if err != nil {
			return nil, err
		}
		header := http.Header{}
		header.Set("Authorization", "Bearer "+jwt)
		modifier := transport.NewHeaderRequestModifier(header)
		return transport.NewTransport(trans, modifier), nil
	}, nil
}

func (a *APIServer) getTargetsForRepoOrError(ctx context.Context, rtBuilder roundTripBuilder, namespace, repoName string) ([]*notaryclient.Target, responses.APIResponse) {
	userConfig, err := a.settingsStore.UserHubConfig()
	if err != nil {
		return nil, responses.APIError(errors.InternalError(ctx, err))
	}

	gun := userConfig.DTRHost + "/" + namespace + "/" + repoName

	rtripper, err := rtBuilder(gun)
	if err != nil {
		return nil, responses.APIError(errors.InternalError(ctx, err))
	}

	err = os.MkdirAll(notaryDataPath, 0600)
	if err != nil {
		return nil, responses.APIError(errors.InternalError(ctx, err))
	}

	var nilFunc notary.PassRetriever
	notaryServerURL, err := url.Parse(fmt.Sprintf("https://%s:%d", containers.NotaryServer.BridgeNameLocalReplica(), deploy.NotaryServerHTTPPort))
	if err != nil {
		return nil, responses.APIError(errors.InternalError(ctx, err))
	}
	nRepo, err := notaryclient.NewNotaryRepository(notaryDataPath, gun, notaryServerURL.String(), rtripper, nilFunc, trustpinning.TrustPinConfig{})
	if err != nil {
		return nil, responses.APIError(errors.InternalError(ctx, err))
	}

	// Retrieve the remote list of signed targets
	targetList := []*notaryclient.Target{}
	targetWithRoleList, err := nRepo.ListTargets()
	switch err.(type) {
	case notaryclient.ErrRepositoryNotExist, notaryclient.ErrRepoNotInitialized:
		if !strings.Contains(userConfig.DTRHost, ":") {
			// if the port is 443, try to check for pushes to `host:443` in addition to `host`, see https://github.com/docker/docker/issues/22644
			gun := userConfig.DTRHost + ":443/" + namespace + "/" + repoName

			rtripper, err := rtBuilder(gun)
			if err != nil {
				return nil, responses.APIError(errors.InternalError(ctx, err))
			}

			nRepo, err := notaryclient.NewNotaryRepository(notaryDataPath, gun, notaryServerURL.String(), rtripper, nilFunc, trustpinning.TrustPinConfig{})
			if err != nil {
				return nil, responses.APIError(errors.InternalError(ctx, err))
			}

			// Retrieve the remote list of signed targets
			targetWithRoleList, err := nRepo.ListTargets()
			switch err.(type) {
			case nil, notaryclient.ErrRepositoryNotExist, notaryclient.ErrRepoNotInitialized:
				break
			default:
				return nil, responses.APIError(errors.InternalError(ctx, err))
			}
			for _, targetWithRole := range targetWithRoleList {
				targetList = append(targetList, &targetWithRole.Target)
			}
		}
		return targetList, nil
	case nil:
		break
	default:
		return nil, responses.APIError(errors.InternalError(ctx, err))
	}
	for _, targetWithRole := range targetWithRoleList {
		targetList = append(targetList, &targetWithRole.Target)
	}
	return targetList, nil
}

// handleGetTagTrust handles checking if a tag has been signed
func (a *APIServer) handleGetTagTrust(ctx context.Context, r *restful.Request) responses.APIResponse {
	namespace := r.PathParameter("namespace")
	repoName := r.PathParameter("reponame")
	tag := r.PathParameter("tag")

	// Gather request data.
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(false),
		makeFilterGetRepoNamespace(namespace),
		makeFilterGetRepository(repoName),
		getRepositoryAccess,
		ensureAccessLevelAtLeastReadOnly,
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	// First check that the tag is in the registry. If not we can bail quickly
	// without querying notary
	schemaTag, err := a.tagMgr.RepositoryTag(namespace+"/"+repoName, tag)
	if err == schema.ErrNoSuchTag {
		return responses.APIError(errors.ErrorCodeNoSuchTag)
	} else if err != nil {
		logrus.Warnf("failed to load tag from DB: %s", err)
		return responses.APIError(errors.InternalError(ctx, err))
	}

	response := responses.MakeTag(schemaTag)

	// Query notary for all targets
	rtBuilder, err := a.createNotaryTransportBuilder(ctx, rd.user, rd.namespace, rd.repo, rd.accessLevel)
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}
	targetList, errResponse := a.getTargetsForRepoOrError(ctx, rtBuilder, namespace, repoName)
	if errResponse != nil {
		return errResponse
	}

	// see if we have the tag in notary
	for _, t := range targetList {
		if t.Name == tag {
			response.InNotary = true
			response.HashMismatch, _ = schemaTag.DigestMismatches(t.Hashes["sha256"])
			break
		}
	}

	return responses.JSONResponse(http.StatusOK, nil, nil, response)
}

func (a *APIServer) deleteTrustDataForRepo(ctx context.Context, rtBuilder roundTripBuilder, namespace, repoName string) error {
	userConfig, err := a.settingsStore.UserHubConfig()
	if err != nil {
		return err
	}

	gun := userConfig.DTRHost + "/" + namespace + "/" + repoName

	rtripper, err := rtBuilder(gun)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(notaryDataPath, 0600); err != nil {
		return err
	}

	var nilFunc notary.PassRetriever
	notaryServerURL, err := url.Parse(fmt.Sprintf("https://%s:%d", containers.NotaryServer.BridgeNameLocalReplica(), deploy.NotaryServerHTTPPort))
	if err != nil {
		return err
	}
	nRepo, err := notaryclient.NewNotaryRepository(notaryDataPath, gun, notaryServerURL.String(), rtripper, nilFunc, trustpinning.TrustPinConfig{})
	if err != nil {
		return err
	}

	// Check if we need to delete the raw reponame or with :443 appended, and then try to delete
	if err := nRepo.Update(false); err != nil {
		switch err.(type) {
		case notaryclient.ErrRepositoryNotExist, notaryclient.ErrRepoNotInitialized:
			if !strings.Contains(userConfig.DTRHost, ":") {
				// if the port is 443, try to check for pushes to `host:443` in addition to `host`, see https://github.com/docker/docker/issues/22644
				expandedGun := userConfig.DTRHost + ":443/" + namespace + "/" + repoName

				rtripper, err := rtBuilder(expandedGun)
				if err != nil {
					return err
				}
				nRepo, err = notaryclient.NewNotaryRepository(notaryDataPath, expandedGun, notaryServerURL.String(), rtripper, nilFunc, trustpinning.TrustPinConfig{})
				if err != nil {
					return err
				}
				if err := nRepo.Update(false); err != nil {
					switch err.(type) {
					case notaryclient.ErrRepositoryNotExist, notaryclient.ErrRepoNotInitialized:
						// If no repo for notary exists, just return nil
						return nil
					default:
						return err
					}
				}
			}
		default:
			return err
		}
	}
	// Delete the remote metadata
	return nRepo.DeleteTrustData(true)
}
