package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/adminserver/api/common/forms"
	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
	"github.com/docker/dhe-deploy/events"
	"github.com/docker/dhe-deploy/manager/schema"
	"github.com/docker/dhe-deploy/registry/client"
	"github.com/docker/dhe-deploy/shared/dtrutil"

	"github.com/docker/distribution/context"
	_ "github.com/docker/distribution/manifest/schema2"

	// drivers for fast repository deletion
	_ "github.com/docker/distribution/registry/storage/driver/azure"
	_ "github.com/docker/distribution/registry/storage/driver/filesystem"
	_ "github.com/docker/distribution/registry/storage/driver/gcs"
	_ "github.com/docker/distribution/registry/storage/driver/inmemory"
	_ "github.com/docker/distribution/registry/storage/driver/middleware/cloudfront"
	_ "github.com/docker/distribution/registry/storage/driver/oss"
	_ "github.com/docker/distribution/registry/storage/driver/s3-aws"
	_ "github.com/docker/distribution/registry/storage/driver/swift"

	enzierrors "github.com/docker/orca/enzi/api/errors"
	"github.com/emicklei/go-restful"
)

// listRepositoriesHandler handles listing all repositories the current user has
// access to. If the client is not authenticated
// then only public repositories are listed.
func (a *APIServer) listRepositoriesHandler(ctx context.Context, r *restful.Request) responses.APIResponse {

	// Gather request data.
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(false),
		getPagerParams,
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	var (
		repos []*schema.Repository
		next  string
		err   error
	)
	if *rd.user.Account.IsAdmin {
		repos, next, err = a.repoMgr.ListAllRepositories(rd.start, rd.limit)
	} else {
		repos, next, err = a.authorizer.SharedRepositoriesForUser(rd.user, rd.start, rd.limit)
		if err != nil {
			return responses.APIError(errors.InternalError(ctx, err))
		}
	}
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}

	accountIDToNameMap := map[string]string{}
	accountIDToTypeMap := map[string]bool{}
	resp := make([]responses.Repository, len(repos))
	for i, repo := range repos {
		namespaceAccountType := accountIDToTypeMap[repo.NamespaceAccountID]
		namespaceAccountName, ok := accountIDToNameMap[repo.NamespaceAccountID]
		if !ok {
			acc, err := rd.user.EnziSession.GetAccount("id:" + repo.NamespaceAccountID)
			apiErrs, _ := err.(*enzierrors.APIErrors)
			if dtrutil.CheckContainsEnziError(apiErrs, "NO_SUCH_ACCOUNT") {
				// TODO delete repository maybe???
				continue
			} else if err != nil {
				return responses.APIError(errors.InternalError(
					ctx,
					fmt.Errorf("received an unexpected error from Enzi: %s", err),
				))
			}

			namespaceAccountName = acc.Name
			namespaceAccountType = acc.IsOrg
			accountIDToNameMap[repo.NamespaceAccountID] = namespaceAccountName
			accountIDToTypeMap[repo.NamespaceAccountID] = namespaceAccountType
		}
		resp[i] = responses.MakeRepository(namespaceAccountName, namespaceAccountType, repo, false)
	}

	return responses.JSONResponsePage(http.StatusOK, nil, nil, responses.Repositories{resp}, r, next, 0)
}

// listNamespaceRepositoriesHandler handles listing all repositories in a given
// namespace that are visible to the client. If the client is not authenticated
// then only public repositories are listed.
func (a *APIServer) listNamespaceRepositoriesHandler(ctx context.Context, r *restful.Request) responses.APIResponse {
	vars := r.PathParameters()
	namespace := vars["namespace"]

	// Gather request data.
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(false),
		makeFilterGetRepoNamespace(namespace),
		getPagerParams,
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	repos, next, err := a.authorizer.VisibleRepositoriesInNamespace(rd.user, rd.namespace, rd.start, rd.limit)
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}

	resp := make([]responses.Repository, len(repos))
	for i, repo := range repos {
		resp[i] = responses.MakeRepository(rd.namespace.Name, rd.namespace.IsOrg, repo, false)
	}

	return responses.JSONResponsePage(http.StatusOK, nil, nil, responses.Repositories{resp}, r, next, 0)
}

// createRepositoryHandler handles creating a new repository in a given
// namespace. The client must be authenticated as a user with "admin" access
// to the repository namespace.
func (a *APIServer) createRepositoryHandler(ctx context.Context, r *restful.Request) responses.APIResponse {
	vars := r.PathParameters()
	namespace := vars["namespace"]

	// Gather request data.
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(true),
		makeFilterGetRepoNamespace(namespace),
		getRepoNamespaceAccess,
		ensureRepoNamespaceAccessLevelAtLeastAdmin,
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	reqOptions := forms.CreateRepo{}
	if err := json.NewDecoder(r.Request.Body).Decode(&reqOptions); err != nil {
		return responses.APIError(errors.MakeError(errors.ErrorCodeInvalidJSON, err))
	}

	if errResponse := validateRepositoryName(namespace, reqOptions.Name); errResponse != nil {
		return errResponse
	}

	if errResponse := validateRepoShortDescription(reqOptions.ShortDescription); errResponse != nil {
		return errResponse
	}

	if reqOptions.Visibility == "" {
		// Use "private" visibility by default.
		reqOptions.Visibility = schema.RepositoryVisibilityPrivate
	} else if errResponse := validateRepoVisibility(reqOptions.Visibility); errResponse != nil {
		return errResponse
	}

	repo := &schema.Repository{
		NamespaceAccountID: rd.namespace.ID,
		Name:               reqOptions.Name,
		ShortDescription:   reqOptions.ShortDescription,
		LongDescription:    reqOptions.LongDescription,
		Visibility:         reqOptions.Visibility,
	}

	if err := a.repoMgr.CreateRepository(repo); err != nil {
		if err == schema.ErrRepositoryExists {
			return responses.APIError(errors.ErrorCodeRepositoryExists)
		}

		return responses.APIError(errors.InternalError(ctx, err))
	}

	if err := events.NewRepositoryEvent(a.eventMgr, rd.user.Account.ID, namespace+"/"+repo.Name); err != nil {
		context.GetLoggerWithField(ctx, "error", err).Error("error creating repo creation event")
	}

	return responses.JSONResponse(http.StatusCreated, nil, nil, responses.MakeRepository(rd.namespace.Name, rd.namespace.IsOrg, repo, true))
}

// getRepositoryHandler handles getting details of a repository. The client
// must be authenticated as a user with at least "read" access to the
// repository OR the repository must have "public" visibility.
func (a *APIServer) getRepositoryHandler(ctx context.Context, r *restful.Request) responses.APIResponse {
	vars := r.PathParameters()
	namespace := vars["namespace"]
	repoName := vars["reponame"]

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

	return responses.JSONResponse(http.StatusOK, nil, nil, responses.MakeRepository(rd.namespace.Name, rd.namespace.IsOrg, rd.repo, true))
}

// patchRepositoryHandler handles updating details of a repository. The client
// must be authenticated as a user with "admin" access to the repository.
func (a *APIServer) patchRepositoryHandler(ctx context.Context, r *restful.Request) responses.APIResponse {
	vars := r.PathParameters()
	namespace := vars["namespace"]
	repoName := vars["reponame"]

	// Gather request data.
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(true),
		makeFilterGetRepoNamespace(namespace),
		makeFilterGetRepository(repoName),
		getRepositoryAccess,
		ensureAccessLevelAtLeastAdmin,
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	reqOptions := forms.UpdateRepo{}
	if err := json.NewDecoder(r.Request.Body).Decode(&reqOptions); err != nil {
		return responses.APIError(errors.MakeError(errors.ErrorCodeInvalidJSON, err))
	}

	if reqOptions.ShortDescription != nil {
		if errResponse := validateRepoShortDescription(*reqOptions.ShortDescription); errResponse != nil {
			return errResponse
		}

		rd.repo.ShortDescription = *reqOptions.ShortDescription
	}

	if reqOptions.LongDescription != nil {
		rd.repo.LongDescription = *reqOptions.LongDescription
	}

	if reqOptions.Visibility != nil {
		if errResponse := validateRepoVisibility(*reqOptions.Visibility); errResponse != nil {
			return errResponse
		}

		rd.repo.Visibility = *reqOptions.Visibility
	}

	// TODO(bbland): allow updating status?

	ruf := schema.RepositoryUpdateFields{
		ShortDescription: &rd.repo.ShortDescription,
		LongDescription:  &rd.repo.LongDescription,
		Visibility:       &rd.repo.Visibility,
	}
	if err := a.repoMgr.UpdateRepository(rd.repo.ID, ruf); err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}
	return responses.JSONResponse(http.StatusOK, nil, nil, responses.MakeRepository(rd.namespace.Name, rd.namespace.IsOrg, rd.repo, true))
}

// deleteRepositoryHandler handles deleting a repository. The client must be
// authenticated as a user with "admin" access to the repository namespace.
//
// This will delete all tags and manifests for the given repository. Layer
// blobs will not be removed until GC runs.
func (a *APIServer) deleteRepositoryHandler(ctx context.Context, r *restful.Request) responses.APIResponse {
	vars := r.PathParameters()
	namespace := vars["namespace"]
	repoName := vars["reponame"]

	// Gather request data.
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(true),
		makeFilterGetRepoNamespace(namespace),
		makeFilterGetRepository(repoName),
		getRepoNamespaceAccess,
		ensureRepoNamespaceAccessLevelAtLeastAdmin,
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	// because we now store state in the blob and metadata store we have a
	// helper client package which handles deletes of repositories
	client, err := client.NewClient(ctx, client.Opts{
		Settings:    a.settingsStore,
		Store:       schema.NewMetadataManager(a.rethinkSession),
		RepoManager: a.repoMgr,
	})

	if err != nil {
		context.GetLoggerWithField(a.baseContext, "error", err).Error("error creating repo client")
		return responses.APIError(errors.InternalError(ctx, err))
	}

	if err = client.DeleteRepository(namespace+"/"+repoName, rd.repo); err != nil {
		context.GetLoggerWithField(a.baseContext, "error", err).Error("error deleting repo")
		return responses.APIError(errors.InternalError(ctx, err))
	}

	// Build a notary roundtripper with admin `*` access to delete
	rtBuilder, err := a.createNotaryTransportBuilder(ctx, rd.user, rd.namespace, rd.repo, rd.accessLevel)
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}
	err = a.deleteTrustDataForRepo(ctx, rtBuilder, namespace, repoName)
	if err != nil {
		context.GetLoggerWithField(ctx, "error", err).Errorf("error deleting repository trust data from notary server: %+v", err)
		return responses.APIError(errors.InternalError(ctx, err))
	}

	if err := events.DeleteRepositoryEvent(a.eventMgr, rd.user.Account.ID, namespace+"/"+repoName); err != nil {
		context.GetLoggerWithField(ctx, "error", err).Error("error creating repo deletion event")
	}

	return responses.JSONResponse(http.StatusNoContent, nil, nil, nil)
}

// getRepositoryTagsHandler handles listing the tags for a repository. The client
// must be authenticated as a user with at least "read" access to the
// repository.
func (a *APIServer) getRepositoryTagsHandler(ctx context.Context, r *restful.Request) responses.APIResponse {
	vars := r.PathParameters()
	namespace := vars["namespace"]
	repoName := vars["reponame"]

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

	var (
		err        error
		tags       []responses.Tag
		notaryTags = map[string][]byte{} // Map of tag name to sha256 checksum
	)

	schemaTags, err := a.tagMgr.RepositoryTags(namespace + "/" + repoName)
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, fmt.Errorf("error listing tags: %v", err)))
	}

	rtBuilder, err := a.createNotaryTransportBuilder(ctx, rd.user, rd.namespace, rd.repo, rd.accessLevel)
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}

	tags = make([]responses.Tag, len(schemaTags))

	// get all signed targets for the given repository from notary
	targetList, respErr := a.getTargetsForRepoOrError(ctx, rtBuilder, namespace, repoName)
	if respErr != nil {
		return respErr
	}
	// build a map from target name to its hash for lookups
	for _, target := range targetList {
		notaryTags[target.Name] = target.Hashes["sha256"]
	}

	// convert schema.Tag into responses.Tag and add notary
	// information to the response struct.
	for i, t := range schemaTags {
		tags[i] = responses.MakeTag(t)
		if dgst, ok := notaryTags[t.Name]; ok {
			tags[i].InNotary = true
			// We can only get an error from DigestMismatches if the registry
			// or notary hash is invalid; in this case someone has manually
			// tampered with data and the tag is invalid.
			//
			// With the guarantees that notary and distribution both provide
			// we can safely swallow this error as digests will always be
			// valid
			tags[i].HashMismatch, _ = t.DigestMismatches(dgst)
		}
	}

	return responses.JSONResponse(http.StatusOK, nil, nil, tags)
}

// getRepositoryManifests handles listing the tags for a repository. The client
// must be authenticated as a user with at least "read" access to the
// repository.
func (a *APIServer) getRepositoryManifestsHandler(ctx context.Context, r *restful.Request) responses.APIResponse {
	vars := r.PathParameters()
	namespace := vars["namespace"]
	repoName := vars["reponame"]

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

	manifests, err := a.mfstMgr.GetRepoManifests(ctx, namespace+"/"+repoName)
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, fmt.Errorf("error listing manifests: %v", err)))
	}

	return responses.JSONResponse(http.StatusOK, nil, nil, responses.MakeManifests(manifests))
}

func (a *APIServer) deleteRepositoryManifestHandler(ctx context.Context, r *restful.Request) responses.APIResponse {
	vars := r.PathParameters()
	namespace := vars["namespace"]
	repoName := vars["reponame"]
	reference := vars["reference"]

	// Gather request data.
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(true),
		makeFilterGetRepoNamespace(namespace),
		makeFilterGetRepository(repoName),
		getRepositoryAccess,
		ensureAccessLevelAtLeastReadWrite,
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	// TODO: Fail unless metadata store is usable
	// because we now store state in the blob and metadata store we have a
	// helper client package which handles deletes of repositories
	client, err := client.NewClient(ctx, client.Opts{
		Settings:    a.settingsStore,
		Store:       schema.NewMetadataManager(a.rethinkSession),
		RepoManager: a.repoMgr,
	})

	if err != nil {
		context.GetLoggerWithField(a.baseContext, "error", err).Error("error creating repo client")
		return responses.APIError(errors.InternalError(ctx, err))
	}

	if err = client.DeleteManifest(namespace+"/"+repoName, reference); err != nil {
		context.GetLoggerWithField(a.baseContext, "error", err).Error("error deleting manifest")
		return responses.APIError(errors.InternalError(ctx, err))
	}

	// TODO: delete manifest event
	// if err = events.DeleteTagEvent(a.eventMgr, rd.user.Account.ID, namespace+"/"+repoName, tag); err != nil {
	// context.GetLoggerWithField(ctx, "error", err).Error("error creating tag deletion event")
	// }

	return responses.JSONResponse(http.StatusNoContent, nil, nil, nil)
}

// deleteRepositoryTagHandler handles deleting a tag for a repository. The
// client must be authenticated as a user with at least "write" access to the
// repository.
//
// If the tag is signed by notary we can't delete the tag; this leaves notary
// and tagstore/blobstore in an incosistent state. Instead we return a 409
// Conflict error.
func (a *APIServer) deleteRepositoryTagHandler(ctx context.Context, r *restful.Request) responses.APIResponse {
	vars := r.PathParameters()
	namespace := vars["namespace"]
	repoName := vars["reponame"]
	tag := vars["tag"]

	// Gather request data.
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(true),
		makeFilterGetRepoNamespace(namespace),
		makeFilterGetRepository(repoName),
		getRepositoryAccess,
		ensureAccessLevelAtLeastReadWrite,
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}
	// TODO: Fail unless metadata store is usable

	// get all signed targets for the given repository from notary
	rtBuilder, err := a.createNotaryTransportBuilder(ctx, rd.user, rd.namespace, rd.repo, rd.accessLevel)
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}
	targetList, respErr := a.getTargetsForRepoOrError(ctx, rtBuilder, namespace, repoName)
	if respErr != nil {
		return respErr
	}
	// see if we have the tag in notary. If so return a 409 response.
	for _, t := range targetList {
		if t.Name == tag {
			return responses.APIError(errors.ErrorCodeTagInNotary)
		}
	}

	// because we now store state in the blob and metadata store we have a
	// helper client package which handles deletes of repositories
	client, err := client.NewClient(ctx, client.Opts{
		Settings:    a.settingsStore,
		Store:       schema.NewMetadataManager(a.rethinkSession),
		RepoManager: a.repoMgr,
	})

	if err != nil {
		context.GetLoggerWithField(a.baseContext, "error", err).Error("error creating repo client")
		return responses.APIError(errors.InternalError(ctx, err))
	}

	if err = client.DeleteTag(namespace+"/"+repoName, tag); err != nil {
		context.GetLoggerWithField(a.baseContext, "error", err).Error("error deleting tag")
		return responses.APIError(errors.InternalError(ctx, err))
	}

	if err = events.DeleteTagEvent(a.eventMgr, rd.user.Account.ID, namespace+"/"+repoName, tag); err != nil {
		context.GetLoggerWithField(ctx, "error", err).Error("error creating tag deletion event")
	}

	return responses.JSONResponse(http.StatusNoContent, nil, nil, nil)
}
