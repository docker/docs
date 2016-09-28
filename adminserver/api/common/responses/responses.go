package responses

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/docker/dhe-deploy/adminserver/api/common"
	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/manager/schema"

	"github.com/docker/distribution/context"
	enziresponses "github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/jose"
	"github.com/emicklei/go-restful"
)

// APIResponse is an interface that is able to write HTTP response headers
// and a body.
type APIResponse interface {
	WriteResponse(ctx context.Context, response *restful.Response)
	ErrorCodes() []string
	StatusCode() int
}

// jsonResponse is an APIResponse that writes a JSON body.
type jsonResponse struct {
	statusCode int
	errorCodes []string
	object     interface{}
	paginated  bool
	next       string
	total      uint
	request    *restful.Request

	header  map[string][]string
	cookies []*http.Cookie
}

// JSONResponse creates an API response with the given status code that writes
// the response as the JSON serialization of the given object.
func JSONResponse(statusCode int, header map[string][]string, cookies []*http.Cookie, object interface{}) APIResponse {
	return &jsonResponse{
		statusCode: statusCode,
		object:     object,
		header:     header,
		cookies:    cookies,
	}
}

func JSONResponsePage(statusCode int, header map[string][]string, cookies []*http.Cookie, object interface{}, r *restful.Request, next string, total uint) APIResponse {
	return &jsonResponse{
		statusCode: statusCode,
		object:     object,
		paginated:  true,
		next:       next,
		total:      total,
		request:    r,
		header:     header,
		cookies:    cookies,
	}
}

func makePageLink(req *http.Request, start string) string {
	values := req.URL.Query()
	values.Set("start", start)
	return fmt.Sprintf("?%s", values.Encode())
}

func (jr *jsonResponse) ErrorCodes() []string {
	return jr.errorCodes
}

// WriteResponse writes this json response to the given http.ResponseWriter.
func (jr *jsonResponse) WriteResponse(ctx context.Context, response *restful.Response) {
	if jr.paginated {
		// Example expected result:
		// Link: </path?q=blah&last=0>; rel="first",
		//  </path?q=blah&last=20>; rel="next"

		limitStr := jr.request.QueryParameter("limit")
		limit, _ := strconv.ParseUint(limitStr, 10, 32)
		if limit == 0 {
			limit = common.DefaultPerPageLimit
		}

		firstLink := makePageLink(jr.request.Request, "")
		links := fmt.Sprintf(`<%s>; rel="first"`, firstLink)
		if jr.next != "" {
			nextLink := makePageLink(jr.request.Request, jr.next)
			links += fmt.Sprintf(`, <%s>; rel="next"`, nextLink)
		}

		response.Header().Add("Link", links)
		if jr.total > 0 {
			response.Header().Add("X-Total-Count", fmt.Sprintf("%d", jr.total))
		}
	}

	if jr.statusCode == http.StatusNoContent {
		response.WriteHeader(jr.statusCode)
		return
	}

	if jr.object == nil {
		jr.object = struct{}{}
	}

	for k, v := range jr.header {
		for _, val := range v {
			response.Header().Add(k, val)
		}
	}

	for _, cookie := range jr.cookies {
		http.SetCookie(response, cookie)
	}

	err := response.WriteHeaderAndEntity(jr.statusCode, jr.object)
	if err != nil {
		context.GetLogger(ctx).Errorf("unable to encode JSON response: %v", err)
	}
}

// StatusCode returns this response's status code.ResponseWriter.
func (jr *jsonResponse) StatusCode() int {
	return jr.statusCode
}

// APIErrorResponse wraps the given apiErrors into a JSON APIResponse with the
// given statusCode.
func APIError(apiErrors ...errors.APIError) APIResponse {
	status := http.StatusOK
	if len(apiErrors) > 0 {
		status = apiErrors[0].HTTPCode
	}
	errorCodes := make([]string, len(apiErrors))
	for i, err := range apiErrors {
		errorCodes[i] = err.Code
	}
	return &jsonResponse{
		statusCode: status,
		errorCodes: errorCodes,
		object: map[string][]errors.APIError{
			"errors": apiErrors,
		},
	}
}

type Team struct {
	ID                 string `json:"id"`
	ClientUserIsMember bool   `json:"clientUserIsMember"`
}

// getTeamResponse returns an interface that cleanly encodes a team in a
// JSON response.
func MakeTeam(teamID string, isMember bool) Team {
	return Team{
		ID:                 teamID,
		ClientUserIsMember: isMember,
	}
}

type Repository struct {
	ID               string `json:"id"`
	Namespace        string `json:"namespace"`
	NamespaceType    string `json:"namespaceType" enum:"user|organization"`
	Name             string `json:"name"`
	ShortDescription string `json:"shortDescription"`
	LongDescription  string `json:"longDescription,omitempty"`
	Visibility       string `json:"visibility" enum:"public|private"`
}

func (r Repository) FullName() string {
	return r.Namespace + "/" + r.Name
}

type Repositories struct {
	Repositories []Repository `json:"repositories"`
}

func MakeRepository(namespaceName string, isOrg bool, repo *schema.Repository, long bool) Repository {
	resp := Repository{
		ID:               repo.ID,
		Namespace:        namespaceName,
		Name:             repo.Name,
		ShortDescription: repo.ShortDescription,
		Visibility:       repo.Visibility,
	}
	if isOrg {
		resp.NamespaceType = "organization"
	} else {
		resp.NamespaceType = "user"
	}

	if long {
		resp.LongDescription = repo.LongDescription
	}

	return resp
}

type Tag struct {
	Name         string    `json:"name"`
	Digest       string    `json:"digest"`
	Author       string    `json:"author"`
	UpdatedAt    time.Time `json:"updatedAt"`
	HashMismatch bool      `json:"hashMismatch" description:"true if the hashes from notary and registry don't match"`
	InNotary     bool      `json:"inNotary" description:"true if the tax exists in Notary"`
	Manifest     Manifest  `json:"manifest"`
}

// MakeTag creates a new responses.Tag given a schema.Tag item created
// from the metadata store.
func MakeTag(t schema.Tag) Tag {
	return Tag{
		Name:      t.Name,
		Digest:    t.Digest,
		Author:    t.Author,
		UpdatedAt: t.UpdatedAt,
		Manifest:  MakeManifest(t.Manifest),
	}
}

type Manifest struct {
	Digest       string    `json:"digest"`
	OS           string    `json:"os"`
	Architecture string    `json:"architecture"`
	MediaType    string    `json:"mediaType"`
	Layers       []string  `json:"layers"`
	Size         int64     `json:"size"`
	CreatedAt    time.Time `json:"createdAt"`
	Author       string    `json:"author"`
}

func MakeManifest(m schema.Manifest) Manifest {
	return Manifest{
		Digest:       m.Digest,
		OS:           m.OS,
		Architecture: m.Architecture,
		MediaType:    m.MediaType,
		Layers:       m.Layers,
		Size:         m.Size,
		CreatedAt:    m.CreatedAt,
		Author:       m.OriginalAuthor,
	}
}

func MakeManifests(list []*schema.Manifest) []Manifest {
	ret := make([]Manifest, len(list))
	for i, m := range list {
		ret[i] = MakeManifest(*m)
	}
	return ret
}

type WithAccessLevel struct {
	AccessLevel string `json:"accessLevel" enum:"read-only|read-write|admin"`
}

type ReplicaSettings struct {
	HTTPPort  uint16 `json:"HTTPPort"`
	HTTPSPort uint16 `json:"HTTPSPort"`
	Node      string `json:"node"`
}

type Settings struct {
	// TODO: document the enum fields
	DTRHost               string                     `json:"dtrHost"`
	ReplicaSettings       map[string]ReplicaSettings `json:"replicaSettings"`
	AuthBypassCA          string                     `json:"authBypassCA"`
	AuthBypassOU          string                     `json:"authBypassOU"`
	HTTPProxy             string                     `json:"httpProxy"`
	HTTPSProxy            string                     `json:"httpsProxy"`
	NoProxy               string                     `json:"noProxy"`
	DisableUpgrades       bool                       `json:"disableUpgrades"`
	ReportAnalytics       bool                       `json:"reportAnalytics"`
	AnonymizeAnalytics    bool                       `json:"anonymizeAnalytics"`
	ReleaseChannel        *string                    `json:"releaseChannel"`
	LogProtocol           string                     `json:"logProtocol"`
	LogHost               string                     `json:"logHost"`
	LogLevel              string                     `json:"logLevel"`
	WebTLSCert            string                     `json:"webTLSCert"`
	WebTLSCA              string                     `json:"webTLSCA"`
	EtcdHeartbeatInterval int                        `json:"etcdHeartbeatInterval"`
	EtcdElectionTimeout   int                        `json:"etcdElectionTimeout"`
	EtcdSnapshotCount     int                        `json:"etcdSnapshotCount"`
	ReplicaID             string                     `json:"replicaID"`
	GCMode                string                     `json:"gcMode"`
}

type UserAccess struct {
	WithAccessLevel
	User enziresponses.Account `json:"user"`
}

type RepoUserAccess struct {
	WithAccessLevel
	User       enziresponses.Account `json:"user"`
	Repository Repository            `json:"repository"`
}

type TeamAccess struct {
	WithAccessLevel
	Team Team `json:"team"`
}

func MakeTeamAccess(rta *schema.RepositoryTeamAccess, isMember bool) TeamAccess {
	return TeamAccess{
		WithAccessLevel: WithAccessLevel{rta.AccessLevel},
		Team:            MakeTeam(rta.TeamID, isMember),
	}
}

type RepoAccess struct {
	WithAccessLevel
	Repository Repository `json:"repository"`
}

func MakeRepoAccess(namespaceName string, isOrg bool, rta *schema.RepositoryTeamAccess) RepoAccess {
	return RepoAccess{
		WithAccessLevel: WithAccessLevel{rta.AccessLevel},
		Repository:      MakeRepository(namespaceName, isOrg, &rta.Repository, false),
	}
}

type RepoTeamAccess struct {
	WithAccessLevel
	Team       Team       `json:"team"`
	Repository Repository `json:"repository"`
}

func MakeRepoTeamAccess(namespaceName string, isOrg bool, rta *schema.RepositoryTeamAccess, isMember bool) RepoTeamAccess {
	return RepoTeamAccess{
		WithAccessLevel: WithAccessLevel{rta.AccessLevel},
		Team:            MakeTeam(rta.TeamID, isMember),
		Repository:      MakeRepository(namespaceName, isOrg, &rta.Repository, false),
	}
}

type Namespace string

func MakeNamespace(ns *enziresponses.Account) Namespace {
	return Namespace(ns.Name)
}

func MakeTeamAccessForNamespace(nta *schema.NamespaceTeamAccess, isMember bool) TeamAccess {
	return TeamAccess{
		WithAccessLevel: WithAccessLevel{nta.AccessLevel},
		Team:            MakeTeam(nta.TeamID, isMember),
	}
}

type NamespaceTeamAccess struct {
	WithAccessLevel
	Team      Team      `json:"team"`
	Namespace Namespace `json:"namespace"`
}

func MakeNamespaceTeamAccess(nta *schema.NamespaceTeamAccess, isMember bool, ns *enziresponses.Account) NamespaceTeamAccess {
	return NamespaceTeamAccess{
		WithAccessLevel: WithAccessLevel{nta.AccessLevel},
		Team:            MakeTeam(nta.TeamID, isMember),
		Namespace:       MakeNamespace(ns),
	}
}

type ListTeamRepoAccess struct {
	Team                 Team         `json:"team"`
	RepositoryAccessList []RepoAccess `json:"repositoryAccessList"`
}

type ListRepoTeamAccess struct {
	Repository     Repository   `json:"repository"`
	TeamAccessList []TeamAccess `json:"teamAccessList"`
}

type ListRepoNamespaceTeamAccess struct {
	Namespace      Namespace    `json:"namespace"`
	TeamAccessList []TeamAccess `json:"teamAccessList"`
}

type DockerRepository struct {
	Description string `json:"description"`
	IsOfficial  bool   `json:"is_official"`
	IsTrusted   bool   `json:"is_trusted"`
	Name        string `json:"name"`
	StarCount   int    `json:"star_count"`
}

type DockerSearch struct {
	NumResults int                `json:"num_results"`
	Query      string             `json:"query"`
	Results    []DockerRepository `json:"results"`
}

// MakeDockerSearch creates a response of the form that `docker search` wants.
func MakeDockerSearch(results []Repository) DockerSearch {
	dockerRepos := make([]DockerRepository, len(results))
	for i, result := range results {
		// We only track name and description. Don't worry about the rest
		dockerRepos[i] = DockerRepository{
			Name:        result.Namespace + "/" + result.Name,
			Description: result.ShortDescription,
		}
	}
	return DockerSearch{
		NumResults: len(results),
		Results:    dockerRepos,
	}
}

type Autocomplete struct {
	AccountResults    []enziresponses.Account `json:"accountResults,omitempty"`
	RepositoryResults []Repository            `json:"repositoryResults,omitempty"`
}

func MakeAutocomplete(accountResults []enziresponses.Account, repositoryResults []Repository) Autocomplete {
	return Autocomplete{
		AccountResults:    accountResults,
		RepositoryResults: repositoryResults,
	}
}

type OpenIDKeys struct {
	Keys []jose.PublicKey `json:"keys"`
}

type ClusterStatus struct {
	RethinkSystemTables map[string]interface{} `json:"rethink_system_tables"`
	EtcdStatus          map[string]interface{} `json:"etcd_status"`
	ReplicaHealth       map[string]string      `json:"replica_health"`
	ReplicaTimestamp    map[string]string      `json:"replica_timestamp"`
	ReplicaRORegistry   map[string]bool        `json:"replica_readonly"`
	GCLockHolder        string                 `json:"gc_lock_holder"`
}

type Events struct {
	Events []schema.Event `json:"events"`
}
