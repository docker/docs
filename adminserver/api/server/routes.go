package server

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"runtime/debug"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/adminserver/api/common/forms"
	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
	"github.com/docker/dhe-deploy/garant/authz"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/jobrunner/jobrunner/jobrunner_config"
	"github.com/docker/dhe-deploy/manager/schema"
	"github.com/docker/dhe-deploy/manager/schema/interfaces"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/api/action_configs"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/api/crons"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/api/jobs"
	jschema "github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema"
	enzierrors "github.com/docker/orca/enzi/api/errors"
	enziresponses "github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/api/server"
	"github.com/docker/orca/enzi/api/server/accounts"
	"github.com/docker/orca/enzi/api/server/admin"
	"github.com/docker/orca/enzi/api/server/config"
	enziJobs "github.com/docker/orca/enzi/api/server/jobs"
	"github.com/docker/orca/enzi/api/server/workers"
	enziSchema "github.com/docker/orca/enzi/schema"

	"github.com/docker/distribution/context"
	"github.com/docker/garant/auth/common"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

// APIServer handles various API requests including managing Teams,
// Repositories, and Repository Access Control.
type APIServer struct {
	baseContext    context.Context
	authorizer     authz.Authorizer
	settingsStore  hubconfig.SettingsStore
	kvStore        hubconfig.KeyValueStore
	propertyMgr    interfaces.PropertyManager
	rethinkSession *rethink.Session
	eventMgr       schema.EventManager
	repoMgr        *schema.RepositoryManager
	tagMgr         *schema.TagManager
	mfstMgr        *schema.ManifestRegistryManager
	metadataMgr    *schema.MetadataManager
	repoAccessMgr  *schema.RepositoryAccessManager
}

// NewAPIServer returns a new APIServer that is ready to use.
func NewAPIServer(baseContext context.Context, authorizer authz.Authorizer, settingsStore hubconfig.SettingsStore, kvStore hubconfig.KeyValueStore, propertyManager interfaces.PropertyManager, session *rethink.Session) *APIServer {
	return &APIServer{
		baseContext:    baseContext,
		authorizer:     authorizer,
		settingsStore:  settingsStore,
		kvStore:        kvStore,
		propertyMgr:    propertyManager,
		rethinkSession: session,
		eventMgr:       schema.NewEventManager(session),
		repoMgr:        schema.NewRepositoryManager(session),
		repoAccessMgr:  schema.NewRepositoryAccessManager(session),
		tagMgr:         schema.NewTagManager(session),
		metadataMgr:    schema.NewMetadataManager(session),
		mfstMgr:        schema.NewManifestRegistryManager(session),
	}
}

// APIHandler is the basic method type for all API request handlers.
type APIHandler func(context.Context, *restful.Request) responses.APIResponse

// wrapHandler wraps the given APIServer handler by setting up the request
// context from this APIServer's baseContext, calling the given handler, and
// writing the response. It also intercepts a panic, logging it with a stack
// traces, writing a JSON error response, and re-panics so that any wrapping
// panic handlers may also handle the panic. Returns a standard http.Handler.
func (a *APIServer) wrapHandler(handler APIHandler) restful.RouteFunction {
	return restful.RouteFunction(func(request *restful.Request, response *restful.Response) {
		ctx := context.WithRequest(a.baseContext, request.Request)
		logger := context.GetRequestLogger(ctx)
		ctx = context.WithLogger(ctx, logger)
		// insert the garant response writer in the middle between restful and its own response writer
		ctx, response.ResponseWriter = context.WithResponseWriter(ctx, response.ResponseWriter)

		defer func() {
			err := recover()
			if err == nil {
				// Not currently panicking.
				return
			}

			// Write a simple error response to the client.
			jsonResponse := responses.APIError(errors.InternalError(ctx, fmt.Errorf("runtime panic: %v", err)))

			jsonResponse.WriteResponse(ctx, response)

			var stack []byte
			if stacker, ok := err.(common.Stacker); ok {
				stack = stacker.Stack()
			} else {
				stack = debug.Stack()
			}

			// Push the stacktrace onto the context to be logged.
			ctx = context.WithValue(ctx, "stackTrace", string(stack))
			ctx = context.WithLogger(ctx, context.GetLogger(ctx, "stackTrace"))
			context.GetResponseLogger(ctx).Errorf("runtime panic: %v", err)

			// Re-panic so that the bugsnag handler (if configured) can notify
			// about the panic as well. The panic will finally be suppressed by
			// the outer-most handler.
			panic(err)
		}()

		apiResponse := handler(ctx, request)
		apiResponse.WriteResponse(ctx, response)

		// log warnings for undocumented routes unless in production
		if !deploy.IsProduction() {
			veMap := validErrors[request.Request.Method+request.SelectedRoutePath()]
			if veMap == nil {
				return
			}

			errorCodes := apiResponse.ErrorCodes()
			if errorCodes == nil {
				return
			}
			if apiResponse.StatusCode() == http.StatusInternalServerError {
				return
			}
			// verify that all error codes returned appear in the list
			// of valid errors
			for _, errorCode := range errorCodes {
				if _, ok := veMap[errorCode]; !ok {
					context.GetLogger(ctx).Warnf("Undocumented error code returned: %s", errorCode)
				}
			}
		}
	})
}

var pathParamRegexp = regexp.MustCompile(`\{([^}]+)\}`)

type pathParam struct {
	Doc  string
	Type string
	Errs []errors.APIError
}

var pathParamDoc = map[string]pathParam{
	"accountname": {
		"user or organization account name",
		"string",
		[]errors.APIError{errors.ErrorCodeNoSuchAccount},
	},
	"username": {"user account name",
		"string",
		[]errors.APIError{errors.ErrorCodeNoSuchUser},
	},
	"orgname": {"organization account name",
		"string",
		[]errors.APIError{errors.ErrorCodeNoSuchOrganization},
	},
	"teamname": {"team name",
		"string",
		[]errors.APIError{errors.ErrorCodeNoSuchTeam},
	},
	"member": {"username of team member",
		"string",
		//XXX: Is this correct?
		[]errors.APIError{errors.ErrorCodeNoSuchUser},
	},
	"namespace": {"namespace/owner of repository",
		"string",
		[]errors.APIError{errors.ErrorCodeNoSuchAccount},
	},
	"reponame": {"name of repository",
		"string",
		[]errors.APIError{errors.ErrorCodeNoSuchRepository},
	},
	"grantee": {"username",
		"string",
		[]errors.APIError{errors.ErrorCodeNoSuchUser},
	},
	"reference": {"digest or tag for an image manifest",
		"string",
		[]errors.APIError{errors.ErrorCodeNoSuchManifest},
	},
	"tag": {"tag name",
		"string",
		[]errors.APIError{errors.ErrorCodeNoSuchTag},
	},
}

type queryParam struct {
	Name         string
	Doc          string
	Type         string
	Required     bool
	DefaultValue string
}

var queryParamDoc = map[string]queryParam{
	"autocompleteQuery": {
		Name:     "query",
		Doc:      "Autocomplete query",
		Type:     "string",
		Required: true,
	},
	"searchQuery": {
		Name:     "query",
		Doc:      "Search query",
		Type:     "string",
		Required: true,
	},
	"start": {
		Doc:          "The ID of the first record on the page",
		Type:         "string",
		DefaultValue: "",
	},
	"limit": {
		Doc:          "Maximum number of results to return",
		Type:         "int",
		DefaultValue: "10",
	},
	"includeRepositories": {
		Doc:          "Whether to include repositories in the response",
		Type:         "boolean",
		DefaultValue: "true",
	},
	"includeAccounts": {
		Doc:          "Whether to include accounts in the response",
		Type:         "boolean",
		DefaultValue: "true",
	},
	"namespace": {
		Doc:  "Exact repository namespace to limit results to.",
		Type: "string",
	},
	"type": {
		Doc:  "Account type to list (either 'user' or 'organization')",
		Type: "string",
	},
}

func (a *APIServer) makeRouteBuilder(ws *restful.WebService, method string, path string, handler APIHandler) *restful.RouteBuilder {
	matches := pathParamRegexp.FindAllSubmatch([]byte(path), -1)

	routeBuilder := ws.Method(method).Path(path).To(a.wrapHandler(handler))

	// Initialize the valid errors struct for this route. This needs to happen before calls to Error
	validErrors[method+ws.RootPath()+path] = map[string]struct{}{}

	for _, match := range matches {
		name := string(match[1])
		param, found := pathParamDoc[name]
		if !found {
			// We should not have undocumented route parameters
			panic(fmt.Sprintf("Undocumented path parameter: %s", name))
		}
		routeBuilder.Param(ws.PathParameter(name, param.Doc).DataType(param.Type))
		routeBuilder.Do(Errors(param.Errs...))
	}
	routeBuilder.Do(Errors(errors.ErrorCodeNotAuthenticated))
	return routeBuilder
}

// XXX: this is global, but it should really be per-application.
// One way to do that is to make Errors me a method on a stateful object
// that can keep track of valid errors
var validErrors = map[string]map[string]struct{}{}

// Errors allows us to use our custom error class to specify possible errors
// Errors must be called only after the route is buildable because it needs
// access to its path
func Errors(errs ...errors.APIError) func(*restful.RouteBuilder) {
	return func(rb *restful.RouteBuilder) {
		for _, err := range errs {
			// register this route
			if !deploy.IsProduction() {
				built := rb.Build()
				id := built.Method + built.Path
				validErrors[id][err.Code] = struct{}{}
			}

			rb.Returns(err.HTTPCode, fmt.Sprintf("%s: %s", err.Code, err.Message), nil)
		}
	}
}

// QueryParams allows us describe query parameters more concisely
func QueryParams(names ...string) func(*restful.RouteBuilder) {
	return func(rb *restful.RouteBuilder) {
		for _, name := range names {
			details, ok := queryParamDoc[name]
			if !ok {
				panic(fmt.Sprintf("Undocumented query parameter: %s", name))
			}
			if details.Name != "" {
				name = details.Name
			}
			rb.Param(restful.QueryParameter(name, details.Doc).
				DataType(details.Type).
				Required(details.Required).
				DefaultValue(details.DefaultValue))
		}
	}
}

// Document that this route is paginated
func Paginated(rb *restful.RouteBuilder) {
	rb.Do(QueryParams("start", "limit"))
}

// WireSubroutes registers API handlers on the given subrouter.
func (a *APIServer) BuildSubroutes(pathPrefix, caller string) (*swagger.Config, *restful.Container) {
	wsContainer := restful.NewContainer()

	accountsWS := new(restful.WebService)
	accountsWS.
		Path(pathPrefix + "/v0/accounts").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("Accounts")

	// Repository Team Access API Endpoints.
	accountsWS.Route(a.makeRouteBuilder(accountsWS, "GET", "/{orgname}/teams/{teamname}/repositoryAccess", a.handleListTeamRepoAccess).
		Operation("ListTeamRepoAccess").
		Doc("List repository access grants for a team").
		Do(Paginated).
		Writes(responses.ListTeamRepoAccess{}).
		Do(Errors(errors.ErrorCodeNotAuthorized)).
		Returns(400, "the repository is not owned by an organization", nil).
		Returns(400, "the team does not belong to the organization", nil).
		Notes(`
*Authorization:* Client must be authenticated as a user who owns the organization the team is in or be a member of that team.
		`))

	// API endpoint to get a user's access level to a repository.
	accountsWS.Route(a.makeRouteBuilder(accountsWS, "GET", "/{username}/repositoryAccess/{namespace}/{reponame}", a.handleGetUserRepoAccess).
		Operation("GetUserRepoAccess").
		Doc("Check a user's access to a repository").
		Writes(responses.RepoUserAccess{}).
		Do(Errors(errors.ErrorCodeNotAuthorized)).
		Notes(`
	*Authorization:* Client must be authenticated either as the user in question or be a system admin.
		`))

	// these endpoints have been disabled for now
	// // API endpoint to get a user's access level to a repository namespace.
	// accountsWS.Route(a.makeRouteBuilder(accountsWS, "GET", "/{username}/repositoryNamespaceAccess/{namespace}", a.handleGetUserRepoNamespaceAccess).
	// 	Operation("GetUserRepoNamespaceAccess").
	// 	Doc("Check a user's access to a repository namespace").
	// 	Writes(responses.RepoNamespaceUserAccess{}).
	// 	Do(Errors(errors.ErrorCodeNotAuthorized)).
	// 	Notes(`
	// *Authorization:* Client must be authenticated either as the user in question or be a system admin.
	// 	`))

	// // API endpoint to list repositories that are explicitly shared with a user.
	// accountsWS.Route(a.makeRouteBuilder(accountsWS, "GET", "/{username}/sharedRepositories", a.handleGetUserSharedRepositories).
	// 	Operation("GetUserSharedRepositories").
	// 	Doc("List repositories that are explicitly shared with a user").
	// 	Writes(responses.Repositories{}).
	// 	Do(Paginated).
	// 	Do(Errors(errors.ErrorCodeNotAuthorized)).
	// 	Notes(`
	// *Authorization:* Client must be authenticated either as the user in question or be a system admin.
	// 	`))

	wsContainer.Add(accountsWS)

	adminWS := new(restful.WebService)
	adminWS.
		Path(pathPrefix + "/v0/meta").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("Admin")

	adminWS.Route(a.makeRouteBuilder(adminWS, "POST", "/settings", a.updateAdminSettingsHandler).
		Operation("UpdateSettings").
		Doc("Update settings").
		Do(Errors(errors.ErrorCodeInvalidSettings)).
		Reads(forms.Settings{}).
		Returns(http.StatusAccepted, "success", nil).
		//Do(Errors(errors.ErrorCodeInvalidSettings)).
		Notes(`
*Authorization:* Client must be authenticated an admin.
		`))

	adminWS.Route(a.makeRouteBuilder(adminWS, "GET", "/settings", a.getAdminSettingsHandler).
		Operation("GetSettings").
		Doc("Get settings").
		Writes(responses.Settings{}).
		Notes(`
*Authorization:* Client must be authenticated an admin.
		`))

	adminWS.Route(a.makeRouteBuilder(adminWS, "GET", "/cluster_status", a.getClusterStatusHandler).
		Operation("GetClusterStatus").
		Doc("Get cluster status").
		Writes(responses.ClusterStatus{}).
		Notes(`
*Authorization:* Client must be authenticated an admin.
		`))

	wsContainer.Add(adminWS)

	reposWS := new(restful.WebService)
	reposWS.
		Path(pathPrefix + "/v0/repositories").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("Repositories")

	// Repository management API endpoints.
	reposWS.Route(a.makeRouteBuilder(reposWS, "GET", "/", a.listRepositoriesHandler).
		Operation("ListRepositories").
		Doc("List all repositories").
		Writes(responses.Repositories{}).
		Do(Paginated).
		Notes(`
*Authorization:* Client must be authenticated as any active user in the system. Results will be filtered to only those repositories visible to the client.
		`))

	reposWS.Route(a.makeRouteBuilder(reposWS, "GET", "/{namespace}", a.listNamespaceRepositoriesHandler).
		Operation("ListNamespaceRepositories").
		Doc("List repositories in a namespace").
		Writes(responses.Repositories{}).
		Do(Paginated).
		Notes(`
*Authorization:* Client must be authenticated as any active user in the system. Results will be filtered to only those repositories visible to the client.
		`))

	reposWS.Route(a.makeRouteBuilder(reposWS, "POST", "/{namespace}", a.createRepositoryHandler).
		Operation("CreateRepository").
		Doc("Create repository").
		Reads(forms.CreateRepo{}).
		Returns(http.StatusCreated, "success", responses.Repository{}).
		Returns(400, "invalid repository details", nil).
		Do(Errors(errors.ErrorCodeNotAuthorized)).
		Do(Errors(errors.ErrorCodeRepositoryExists)).
		Notes(`
*Authorization:* Client must be authenticated as a user who has admin access to the
repository namespace (i.e., user owns the repo or is a member of a team with
"admin" level access to the organization's namespace of repositories).
		`))

	reposWS.Route(a.makeRouteBuilder(reposWS, "GET", "/{namespace}/{reponame}", a.getRepositoryHandler).
		Operation("GetRepository").
		Doc("View details of a repository").
		Writes(responses.Repository{}).
		Notes(`
*Authorization:* Client must be authenticated as a user who has visibility to the repository.
		`))

	reposWS.Route(a.makeRouteBuilder(reposWS, "PATCH", "/{namespace}/{reponame}", a.patchRepositoryHandler).
		Operation("PatchRepository").
		Doc("Update details of a repository").
		Reads(forms.UpdateRepo{}).
		Writes(responses.Repository{}).
		Returns(400, "invalid repository details", nil).
		Do(Errors(errors.ErrorCodeNotAuthorized)).
		Do(Errors(errors.ErrorCodeInvalidRepositoryShortDescription)).
		Do(Errors(errors.ErrorCodeInvalidRepositoryVisibility)).
		Notes(`
*Authorization:* Client must be authenticated as a user who has "admin" access to the repository
(i.e., user owns the repo or is a member of a team with "admin" level access to the organization"s repository).

Note that a repository cannot be renamed this way.
		`))

	reposWS.Route(a.makeRouteBuilder(reposWS, "DELETE", "/{namespace}/{reponame}", a.deleteRepositoryHandler).
		Operation("DeleteRepository").
		Doc("Remove a repository").
		Returns(http.StatusNoContent, "success or repository does not exist", nil).
		Do(Errors(errors.ErrorCodeNotAuthorized)).
		Notes(`
*Authorization:* Client must be authenticated as a user who has "admin" access to the repository
(i.e., user owns the repo or is a member of a team with "admin" level access to the organization"s repository).
		`))

	// Repository Team Access API Endpoints.
	reposWS.Route(a.makeRouteBuilder(reposWS, "GET", "/{namespace}/{reponame}/teamAccess", a.handleListRepoTeamAccess).
		Operation("ListRepoTeamAccess").
		Doc("List teams granted access to an organization-owned repository").
		Do(Paginated).
		Writes(responses.ListRepoTeamAccess{}).
		Returns(400, "the repository is not owned by an organization", nil).
		Do(Errors(errors.ErrorCodeNotAuthorized)).
		Notes(`
*Authorization:* Client must be authenticated as a user who has "admin" level access to the repository.
		`))

	reposWS.Route(a.makeRouteBuilder(reposWS, "PUT", "/{namespace}/{reponame}/teamAccess/{teamname}", a.handleGrantRepoTeamAccess).
		Operation("GrantRepoTeamAccess").
		Doc("Set a team's access to an orgnization-owned repository").
		Reads(forms.Access{}).
		Writes(responses.RepoTeamAccess{}).
		Returns(400, "the repository is not owned by an organization", nil).
		Returns(400, "the team does not belong to the organization", nil).
		Do(Errors(errors.ErrorCodeNotAuthorized)).
		Notes(`
*Authorization:* Client must be authenticated as a user who has "admin" level access to the repository.
		`))

	reposWS.Route(a.makeRouteBuilder(reposWS, "DELETE", "/{namespace}/{reponame}/teamAccess/{teamname}", a.handleRevokeRepoTeamAccess).
		Operation("RevokeRepoTeamAccess").
		Doc("Revoke a team's acccess to an organization-owned repository").
		Returns(http.StatusNoContent, "success or the team is not in the access list or there is no such team in the organization", nil).
		Returns(400, "the repository is not owned by an organization", nil).
		Do(Errors(errors.ErrorCodeNotAuthorized)).
		Notes(`
*Authorization:* Client must be authenticated as a user who has "admin" level access to the repository.
		`))

	reposWS.Route(a.makeRouteBuilder(reposWS, "GET", "/{namespace}/{reponame}/tags", a.getRepositoryTagsHandler).
		Operation("ListRepoTags").
		Doc("List the available tags for a repository").
		Writes([]responses.Tag{}).
		Notes(`
*Authorization:* Client must be authenticated as a user who has visibility to the repository.
		`))

	reposWS.Route(a.makeRouteBuilder(reposWS, "GET", "/{namespace}/{reponame}/tags/{tag}/trust", a.handleGetTagTrust).
		Operation("GetTrustForTag").
		Doc("Get Notary trust info about a specific tag").
		Writes(responses.Tag{}).
		Notes(`
Repository results will be filtered to only those repositories visible to the client.
		`))

	reposWS.Route(a.makeRouteBuilder(reposWS, "DELETE", "/{namespace}/{reponame}/tags/{tag}", a.deleteRepositoryTagHandler).
		Operation("DeleteRepoTag").
		Doc("Delete a tag for a repository").
		Do(Errors(errors.ErrorCodeNotAuthenticated, errors.ErrorCodeNotAuthorized, errors.ErrorCodeTagInNotary)).
		Returns(http.StatusNoContent, "success", nil).
		Notes(`
*Authorization:* Client must be authenticated as a user who has "write" level access to the repository.
		`))

	// Manifests
	reposWS.Route(a.makeRouteBuilder(reposWS, "GET", "/{namespace}/{reponame}/manifests", a.getRepositoryManifestsHandler).
		Operation("ListRepoManifests").
		Doc("List the available manifests for a repository").
		Writes([]responses.Manifest{}).
		Notes(`
*Authorization:* Client must be authenticated as a user who has visibility to the repository.
		`))
	reposWS.Route(a.makeRouteBuilder(reposWS, "DELETE", "/{namespace}/{reponame}/manifests/{reference}", a.deleteRepositoryManifestHandler).
		Operation("DeleteRepoManifest").
		Doc("Delete a manifest for a repository").
		Do(Errors(errors.ErrorCodeNotAuthenticated, errors.ErrorCodeNotAuthorized)).
		Returns(http.StatusOK, "success", nil).
		Notes(`
*Authorization:* Client must be authenticated as a user who has "write" level access to the repository.
		`))

	wsContainer.Add(reposWS)

	reposNSWS := new(restful.WebService)
	reposNSWS.
		Path(pathPrefix + "/v0/repositoryNamespaces").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("Repository Namespaces")

	// Repository Namespace Team Access API Endpoints.
	reposNSWS.Route(a.makeRouteBuilder(reposNSWS, "GET", "/{namespace}/teamAccess", a.handleListRepoNamespaceTeamAccess).
		Operation("ListRepoNamespaceTeamAccess").
		Doc("List teams granted access to an organization-owned namespace of repositories").
		Do(Paginated).
		Writes(responses.ListRepoNamespaceTeamAccess{}).
		Returns(400, "the namespace is not owned by an organization", nil).
		Do(Errors(errors.ErrorCodeNotAuthorized)).
		Notes(`
*Authorization:* Client must be authenticated as a user who has ‘admin’ level access to the namespace.
		`))

	reposNSWS.Route(a.makeRouteBuilder(reposNSWS, "GET", "/{namespace}/teamAccess/{teamname}", a.handleGetRepoNamespaceTeamAccess).
		Operation("GetRepoNamespaceTeamAccess").
		Doc("Get a team's granted access to an organization-owned namespace of repositories").
		Writes(responses.NamespaceTeamAccess{}).
		Do(Errors(errors.ErrorCodeNotAuthorized)).
		Notes(`
*Authorization:* Client must be authenticated as a user who has "admin" level access to
the namespace, is a system admin, member of the organization's "owners" team, or is a
member of the team in question.
		`))

	reposNSWS.Route(a.makeRouteBuilder(reposNSWS, "PUT", "/{namespace}/teamAccess/{teamname}", a.handleGrantRepoNamespaceTeamAccess).
		Operation("GrantRepoNamespaceTeamAccess").
		Doc("Set a team's access to an organization-owned namespace of repositories").
		Reads(forms.Access{}).
		Writes(responses.NamespaceTeamAccess{}).
		Returns(400, "the namespace is not owned by an organization", nil).
		Returns(400, "the team does not belong to the owning organization", nil).
		Do(Errors(errors.ErrorCodeNotAuthorized)).
		Notes(`
*Authorization:* Client must be authenticated as a user who has ‘admin’ level access to the namespace.
		`))

	reposNSWS.Route(a.makeRouteBuilder(reposNSWS, "DELETE", "/{namespace}/teamAccess/{teamname}", a.handleRevokeRepoNamespaceTeamAccess).
		Operation("RevokeRepoNamespaceTeamAccess").
		Doc("Revoke a team's access to an organization-owned namespace of repositories").
		Returns(http.StatusNoContent, "success or the team does not exist in the access list or there is no such team in the organization", nil).
		Do(Errors(errors.ErrorCodeNotAuthorized)).
		Notes(`
*Authorization:* Client must be authenticated as a user who has ‘admin’ level access to the namespace.
		`))

	wsContainer.Add(reposNSWS)

	indexWS := new(restful.WebService)
	indexWS.
		Path(pathPrefix + "/v0/index").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("Index")

	indexWS.Route(a.makeRouteBuilder(indexWS, "GET", "/dockersearch", a.handleDockerSearch).
		Operation("Docker Search").
		Doc("Search Docker repositories").
		Do(QueryParams("searchQuery")).
		Writes(responses.DockerSearch{}).
		Notes(`
This is used for the Docker CLI's docker search command. Repository results will be filtered to only those repositories visible to the client.
		`))

	indexWS.Route(a.makeRouteBuilder(indexWS, "GET", "/autocomplete", a.handleAutocomplete).
		Operation("Autocomplete").
		Doc("Autocompletion for repositories and/or accounts").
		Do(QueryParams("autocompleteQuery")).
		Do(QueryParams("includeRepositories")).
		Do(QueryParams("includeAccounts")).
		Do(QueryParams("namespace")).
		Writes(responses.Autocomplete{}).
		Notes(`
Repository results will be filtered to only those repositories visible to the client. Account results will not be filtered.
		`))

	wsContainer.Add(indexWS)

	eventsWS := new(restful.WebService)
	eventsWS.
		Path(pathPrefix + "/v0/events").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("Events")

	eventsWS.Route(a.makeRouteBuilder(eventsWS, "GET", "/", a.getEvents).
		Operation("GetEvents").
		Doc("Get Events").
		Writes(responses.Events{}))

	wsContainer.Add(eventsWS)

	// documentation-only services
	if caller == "apidocgen" {
		enziServices := getEnziServices("enzi")
		for _, webService := range enziServices {
			wsContainer.Add(webService)
		}
	}

	// shared services with enzi
	schemaMgr := jschema.NewJobrunnerManager(deploy.JobrunnerDBName, a.rethinkSession)
	jobsService, err := jobs.NewService(context.Background(), schemaMgr, jobrunner_config.RegisteredActions, path.Join(pathPrefix+"/v0/jobs"), makeAuthWrapper(a))
	if err != nil {
		log.Fatalf("unable to create jobs service: %s", err)
	}
	wsContainer.Add(jobsService.WebService)

	cronsService, err := crons.NewService(context.Background(), schemaMgr, jobrunner_config.RegisteredActions, path.Join(pathPrefix+"/v0/crons"), makeAuthWrapper(a))
	if err != nil {
		log.Fatalf("unable to create crons service: %s", err)
	}
	wsContainer.Add(cronsService.WebService)

	actionConfigsService, err := action_configs.NewService(context.Background(), schemaMgr, jobrunner_config.RegisteredActions, path.Join(pathPrefix+"/v0/action_configs"), makeAuthWrapper(a))
	if err != nil {
		log.Fatalf("unable to create actionConfigs service: %s", err)
	}
	wsContainer.Add(actionConfigsService.WebService)

	// if this fails, it's a packaging error, but it's not critical so we ignore it if it ever happens
	description, _ := ioutil.ReadFile("/ui/swagger/api_intro.md")

	config := swagger.Config{
		WebServices:     wsContainer.RegisteredWebServices(),
		ApiPath:         pathPrefix + "/docs.json",
		SwaggerPath:     pathPrefix + "/docs/",
		SwaggerFilePath: "/swagger",
		Info: swagger.Info{
			Title:       fmt.Sprintf("DTR %s API Documentation", deploy.ShortVersion),
			Description: string(description),
		},
	}

	swagger.RegisterSwaggerService(config, wsContainer)

	// We add the openidWS after Registering web services with swagger since we don't want it to be documented
	openIDWS := new(restful.WebService)
	openIDWS.
		Path(pathPrefix + "/v0/openid").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	openIDWS.Route(a.makeRouteBuilder(openIDWS, "GET", "/begin", a.handleOpenIDBegin).
		Operation("OpenIDBegin"))

	openIDWS.Route(a.makeRouteBuilder(openIDWS, "GET", "/callback", a.handleOpenIDCallback).
		Operation("OpenIDCallback"))

	openIDWS.Route(a.makeRouteBuilder(openIDWS, "GET", "/keys", a.handleOpenIDKeys).
		Operation("OpenIDKeys").
		Writes(responses.OpenIDKeys{}))

	wsContainer.Add(openIDWS)

	return &config, wsContainer
}
func makeAuthWrapper(a *APIServer) func(server.Handler) server.Handler {
	return func(handler server.Handler) server.Handler {
		// special enzi authentication wrapper
		return server.Handler(func(context context.Context, req *restful.Request) enziresponses.APIResponse {
			rd := newRequestData(a, context, req)
			if rd.addFilters(
				makeFilterGetAuthenticatedUser(true),
				ensureIsAdmin,
			).evaluateFilters(); rd.errResponse != nil {
				// XXX: this masks the error, but it should be fine...
				return enziresponses.APIError(enzierrors.NotAuthorized("system admin access is required"))
			}
			return handler(context, req)
		})
	}
}

func getEnziServices(pathPrefix string) []*restful.WebService {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	schemaMgr := enziSchema.NewRethinkDBManager(nil)
	workerClient := server.NewHTTPClient(tlsConfig)

	accountsService := accounts.NewService(context.Background(), schemaMgr, path.Join("/", pathPrefix, "v0/accounts"))
	adminService := admin.NewService(context.Background(), schemaMgr, path.Join("/", pathPrefix, "v0/admin"))
	configService := config.NewService(context.Background(), schemaMgr, path.Join("/", pathPrefix, "v0/config"))
	workersService := workers.NewService(context.Background(), schemaMgr, workerClient, path.Join("/", pathPrefix, "v0/workers"))
	jobsService, err := enziJobs.NewService(context.Background(), schemaMgr, workerClient, path.Join("/", pathPrefix, "v0/jobs"))
	if err != nil {
		log.Fatalf("unable to create jobs service: %s", err)
	}
	return []*restful.WebService{accountsService.WebService, adminService.WebService, configService.WebService, workersService.WebService, jobsService.WebService}
}
