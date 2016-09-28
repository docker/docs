package server

import (
	"net/http"

	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
	"github.com/docker/dhe-deploy/manager/schema"
	"github.com/docker/distribution/context"
	"github.com/emicklei/go-restful"
)

// handleGetUserRepoAccess queries the level of access that the specified user
// has on the specified repository. This is done using all access grants that
// the user may have on the repository either through org/team membership or
// grants by the repository owner or admin if it is a user-owned repository.
// The client must be authenticated as a system admin or the user in question.
func (a *APIServer) handleGetUserRepoAccess(ctx context.Context, r *restful.Request) responses.APIResponse {
	vars := r.PathParameters()
	targetUsername := vars["username"]
	namespace := vars["namespace"]
	repoName := vars["reponame"]

	// Gather request data.
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(true),
		makeFilterGetRepoNamespace(namespace),
		makeFilterGetRepository(repoName),
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	if targetUsername != rd.user.Account.Name {
		return responses.APIError(errors.NotAuthorizedError("You may only request your own repository user access"))
	}

	accessLevel, err := a.authorizer.RepositoryAccess(rd.user, rd.repo, rd.namespace)
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}

	if !(*rd.user.Account.IsAdmin || schema.AccessLevelAtLeast(accessLevel, schema.AccessLevelReadOnly)) {
		// This repository is not visible to the user. Return "404 Not Found".
		return responses.APIError(errors.NoSuchRepositoryError(namespace, repoName))
	}

	rua := &responses.RepoUserAccess{
		WithAccessLevel: responses.WithAccessLevel{accessLevel},
		Repository:      responses.MakeRepository(rd.namespace.Name, rd.namespace.IsOrg, rd.repo, false),
		User:            *rd.user.Account,
	}

	return responses.JSONResponse(http.StatusOK, nil, nil, rua)
}
