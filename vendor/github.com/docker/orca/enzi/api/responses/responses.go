package responses

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/emicklei/go-restful"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api"
	"github.com/docker/orca/enzi/api/errors"
	ldapconfig "github.com/docker/orca/enzi/authn/ldap/config"
	"github.com/docker/orca/enzi/config"
	"github.com/docker/orca/enzi/schema"
)

// APIResponse is an interface that is able to write HTTP response headers
// and a body.
type APIResponse interface {
	WriteResponse(ctx context.Context, response *restful.Response)
	AddCookies(cookies ...*http.Cookie)
	StatusCode() int
}

// jsonResponse is an APIResponse that writes a JSON body.
type jsonResponse struct {
	statusCode    int
	cookiesToSet  []*http.Cookie
	object        interface{}
	paginated     bool
	nextPageStart string
	request       *restful.Request
	header        http.Header
}

// JSONResponse creates an API response with the given status code that writes
// the response as the JSON serialization of the given object.
func JSONResponse(statusCode int, object interface{}) APIResponse {
	return JSONResponseWithHeaders(statusCode, object, nil)
}

// JSONResponseWithHeaders creates an API response with the given status code
// that writes the response as the JSON serialization of the given object and
// adds additiona headers from the given header value.
func JSONResponseWithHeaders(statusCode int, object interface{}, header http.Header) APIResponse {
	return &jsonResponse{
		statusCode: statusCode,
		object:     object,
		header:     header,
	}
}

// JSONResponsePage creates a JSON repsponse with additional response header
// values for paginated results.
func JSONResponsePage(statusCode int, object interface{}, r *restful.Request, nextPageStart string) APIResponse {
	return &jsonResponse{
		statusCode:    statusCode,
		object:        object,
		paginated:     true,
		nextPageStart: nextPageStart,
		request:       r,
	}
}

func makePageLink(req *http.Request, pageStart string) string {
	values := req.URL.Query()
	values.Set("start", pageStart)
	return fmt.Sprintf("?%s", values.Encode())
}

// WriteResponse writes this json response to the given http.ResponseWriter.
func (jr *jsonResponse) WriteResponse(ctx context.Context, response *restful.Response) {
	for header, values := range jr.header {
		for _, value := range values {
			response.AddHeader(header, value)
		}
	}

	for _, cookie := range jr.cookiesToSet {
		http.SetCookie(response.ResponseWriter, cookie)
	}

	if jr.paginated {
		// Example expected result:
		// Link: </path?q=blah&last=0>; rel="first",
		//  </path?q=blah&last=20>; rel="next"

		limitStr := jr.request.QueryParameter("limit")
		limit, _ := strconv.ParseUint(limitStr, 10, 32)
		if limit == 0 {
			limit = api.DefaultPerPageLimit
		}

		firstLink := makePageLink(jr.request.Request, "")
		links := fmt.Sprintf(`<%s>; rel="first"`, firstLink)
		if jr.nextPageStart != "" {
			nextLink := makePageLink(jr.request.Request, jr.nextPageStart)
			links += fmt.Sprintf(`, <%s>; rel="next"`, nextLink)

			// For simplicity.
			response.Header().Set("X-Next-Page-Start", jr.nextPageStart)
		}

		response.Header().Add("Link", links)
	}

	if jr.statusCode == http.StatusNoContent {
		response.WriteHeader(jr.statusCode)
		return
	}

	if jr.object == nil {
		jr.object = struct{}{}
	}

	err := response.WriteHeaderAndEntity(jr.statusCode, jr.object)
	if err != nil {
		context.GetLogger(ctx).Errorf("unable to encode JSON response: %v", err)
	}
}

func (jr *jsonResponse) AddCookies(cookies ...*http.Cookie) {
	jr.cookiesToSet = append(jr.cookiesToSet, cookies...)
}

// StatusCode returns this response's status code.ResponseWriter.
func (jr *jsonResponse) StatusCode() int {
	return jr.statusCode
}

// APIError wraps the given apiErrors into a JSON APIResponse with the
// given statusCode.
func APIError(apiErrors ...*errors.APIError) APIResponse {
	status := http.StatusInternalServerError
	if len(apiErrors) > 0 {
		status = apiErrors[0].HTTPCode
	}

	return &jsonResponse{
		statusCode: status,
		object: errors.APIErrors{
			Errors: apiErrors,
		},
	}
}

// An Account object is used as a user or organization account response.
type Account struct {
	Name     string `json:"name"     description:"Name of the account"`
	ID       string `json:"id"       description:"ID of the account"`
	FullName string `json:"fullName" description:"Full Name of the account"`
	IsOrg    bool   `json:"isOrg"    description:"Whether the account is an organization (or user)"`

	// Fields for users only.
	IsAdmin  *bool `json:"isAdmin,omitempty"  description:"Whether the user is a system admin (users only)"`
	IsActive *bool `json:"isActive,omitempty" description:"Whether the user is active and can login (users only)"`
}

// MakeUser converts a user from the schema backend to an account response
// object.
func MakeUser(user *schema.Account) Account {
	return Account{
		Name:     user.Name,
		ID:       user.ID,
		IsOrg:    false,
		FullName: user.FullName,
		IsAdmin:  &user.IsAdmin,
		IsActive: &user.IsActive,
	}
}

// MakeOrg converts an organization from the schema backend to an account
// response object.
func MakeOrg(org *schema.Account) Account {
	return Account{
		Name:     org.Name,
		ID:       org.ID,
		IsOrg:    true,
		FullName: org.FullName,
	}
}

// MakeAccount converts a user or organizaiton from the schema backend to an
// account response object.
func MakeAccount(acct *schema.Account) Account {
	if acct.IsOrg {
		return MakeOrg(acct)
	}

	// Otherwise it must be a user type account.
	return MakeUser(acct)
}

// An Accounts object is used as a response for a list of accounts.
type Accounts struct {
	Accounts []Account `json:"accounts"`
}

// MakeAccounts converts a slice of users and/or organizaitons from the schema
// backend to an accounts list response object.
func MakeAccounts(accts []schema.Account) Accounts {
	accounts := make([]Account, len(accts))
	for i := range accts {
		accounts[i] = MakeAccount(&accts[i])
	}

	return Accounts{accounts}

}

// A Team object is used as a team response.
type Team struct {
	OrgID       string `json:"orgID"                    description:"ID of the organizaiton to which this team belongs"`
	Name        string `json:"name"                     description:"Name of the team"`
	ID          string `json:"id"                       description:"ID of the team"`
	Description string `json:"description"              description:"Description of the team"`
}

// MakeTeam converts a team from the schema backend to a team response object.
func MakeTeam(team *schema.Team) Team {
	return Team{
		OrgID:       team.OrgID,
		Name:        team.Name,
		ID:          team.ID,
		Description: team.Description,
	}
}

// A Teams object is used as a response for a list of teams.
type Teams struct {
	Teams []Team `json:"teams"`
}

// MakeTeams converts a slice of teams from the schema backend to a teams list
// response object.
func MakeTeams(teams []schema.Team) Teams {
	teamObjects := make([]Team, len(teams))
	for i := range teams {
		teamObjects[i] = MakeTeam(&teams[i])
	}

	return Teams{teamObjects}

}

// A Member object is used as a response for information about a given user's
// membership in some organization or team.
type Member struct {
	Member   Account `json:"member"   description:"The user which is a member of the organization or team"`
	IsAdmin  bool    `json:"isAdmin"  description:"Whether the member is an admin of the organization or team"`
	IsPublic bool    `json:"isPublic" description:"Whether the membership is public"`
}

// MakeMember conversts a MemberInfo from the schema backend to a member
// response object.
func MakeMember(membership *schema.MemberInfo) Member {
	return Member{
		Member:   MakeUser(&membership.Member),
		IsAdmin:  membership.IsAdmin,
		IsPublic: membership.IsPublic,
	}
}

// A Members object is used as a response for a list of members.
type Members struct {
	Members []Member `json:"members"`
}

// MakeMembers converts a slice of MemberInfo from the schema backend to a
// Members list response object.
func MakeMembers(memberInfos []schema.MemberInfo) Members {
	members := make([]Member, len(memberInfos))
	for i := range memberInfos {
		members[i] = MakeMember(&memberInfos[i])
	}

	return Members{members}
}

// A MemberOrg object is used as a response for information about some user's
// membership in a given organization.
type MemberOrg struct {
	Org      Account `json:"org"      description:"The organization which the user is a member of"`
	IsAdmin  bool    `json:"isAdmin"  description:"Whether the user is an admin of the organization"`
	IsPublic bool    `json:"isPublic" description:"Whether the user is a public member of the organization"`
}

// MakeMemberOrg converts a memberOrg from the schema backend to a memberOrg
// response object.
func MakeMemberOrg(memberOrg *schema.MemberOrg) MemberOrg {
	return MemberOrg{
		Org:      MakeOrg(&memberOrg.Org),
		IsAdmin:  memberOrg.IsAdmin,
		IsPublic: memberOrg.IsPublic,
	}
}

// A MemberOrgs object is used as a response for a list of memberOrgs.
type MemberOrgs struct {
	MemberOrgs []MemberOrg `json:"memberOrgs"`
}

// MakeMemberOrgs converts a slice of MemberOrg from the schema backend to a
// MemberOrgs list response object.
func MakeMemberOrgs(memberOrgs []schema.MemberOrg) MemberOrgs {
	memberOrgObjects := make([]MemberOrg, len(memberOrgs))
	for i := range memberOrgs {
		memberOrgObjects[i] = MakeMemberOrg(&memberOrgs[i])
	}

	return MemberOrgs{memberOrgObjects}
}

// A MemberTeam object is used as a response for information about some user's
// membership in a given team.
type MemberTeam struct {
	Team     Team `json:"team"     description:"The team which the user is a member of"`
	IsAdmin  bool `json:"isAdmin"  description:"Whether the user is an admin of the team"`
	IsPublic bool `json:"isPublic" description:"Whether the user is a public member of the team"`
}

// MakeMemberTeam converts a memberTeam from the schema backend to a memberTeam
// response object.
func MakeMemberTeam(memberTeam *schema.MemberTeam) MemberTeam {
	return MemberTeam{
		Team:     MakeTeam(&memberTeam.Team),
		IsAdmin:  memberTeam.IsAdmin,
		IsPublic: memberTeam.IsPublic,
	}
}

// A MemberTeams object is used as a response for a list of memberTeams.
type MemberTeams struct {
	MemberTeams []MemberTeam `json:"memberTeams"`
}

// MakeMemberTeams converts a slice of MemberTeam from the schema backend to a
// MemberTeams list response object.
func MakeMemberTeams(memberTeams []schema.MemberTeam) MemberTeams {
	memberTeamObjects := make([]MemberTeam, len(memberTeams))
	for i := range memberTeams {
		memberTeamObjects[i] = MakeMemberTeam(&memberTeams[i])
	}

	return MemberTeams{memberTeamObjects}
}

// A MemberSyncOpts object contains options for syncing members of an
// organization or team.
type MemberSyncOpts struct {
	EnableSync         bool   `json:"enableSync"         description:"Whether to enable LDAP syncing. If false, all other fields are ignored"`
	SelectGroupMembers bool   `json:"selectGroupMembers" description:"Whether to sync using a group DN and member attribute selection or to use a search filter (if false)"`
	GroupDN            string `json:"groupDN"            description:"The distinguished name of the LDAP group. Applicable only if selectGroupMembers is true, ignored otherwise"`
	GroupMemberAttr    string `json:"groupMemberAttr"    description:"The name of the LDAP group entry attribute which corresponds to distinguished names of members. Applicable only if selectGroupMembers is true, ignored otherwise"`
	SearchBaseDN       string `json:"searchBaseDN"       description:"The distinguished name of the element from which the LDAP server will search for users. Applicable only if selectGroupMembers is false, ignored otherwise"`
	SearchScopeSubtree bool   `json:"searchScopeSubtree" description:"Whether to search for users in the entire subtree of the base DN or to only search one level under the base DN (if false). Applicable only if selectGroupMembers is false, ignored otherwise"`
	SearchFilter       string `json:"searchFilter"       description:"The LDAP search filter used to select users if selectGroupMembers is false, may be left blank"`
}

// MakeMemberSyncOpts returns a MemberSyncOpts response object for the given
// sync options from the schema backend.
func MakeMemberSyncOpts(opts schema.MemberSyncOpts) MemberSyncOpts {
	return MemberSyncOpts{
		EnableSync:         opts.EnableSync,
		SelectGroupMembers: opts.SelectGroupMembers,
		GroupDN:            opts.GroupDN,
		GroupMemberAttr:    opts.GroupMemberAttr,
		SearchBaseDN:       opts.SearchBaseDN,
		SearchScopeSubtree: opts.SearchScopeSubtree,
		SearchFilter:       opts.SearchFilter,
	}
}

// An AuthConfig object is used as a response for system auth configuration.
type AuthConfig struct {
	Backend string `json:"backend" description:"The name of the auth backend to use" enum:"managed|ldap"`
}

// MakeAuthConfig returns an AuthConfig response object for the given config.
func MakeAuthConfig(authConfig *config.Auth) AuthConfig {
	return AuthConfig{
		Backend: authConfig.Backend,
	}
}

// LDAPSettings is used in the AuthConfig response.
type LDAPSettings struct {
	RecoveryAdminUsername string           `json:"recoveryAdminUsername" description:"The user with this name will be able to authenticate to the system even when the LDAP server is unavailable or if auth is misconfigured"`
	ServerURL             string           `json:"serverURL"             description:"The URL of the LDAP server"`
	NoSimplePagination    bool             `json:"noSimplePagination"    description:"The server does not support the Simple Paged Results control extension (RFC 2696)"`
	StartTLS              bool             `json:"startTLS"              description:"Whether to use StartTLS to secure the connection to the server, ignored if server URL scheme is 'ldaps://'"`
	RootCerts             string           `json:"rootCerts"             description:"A root certificate bundle to use when establishing a TLS connection to the server"`
	TLSSkipVerify         bool             `json:"tlsSkipVerify"         description:"Whether to skip verifying of the server's certificate when establishing a TLS connection, not recommended unless testing on a secure network"`
	ReaderDN              string           `json:"readerDN"              description:"The distinguished name the system will use to bind to the LDAP server when performing searches"`
	ReaderPassword        string           `json:"readerPassword"        description:"The password that the system will use to bind to the LDAP server when performing searches"`
	UserSearchConfigs     []UserSearchOpts `json:"userSearchConfigs"     description:"One or more settings for syncing users"`
	AdminSyncOpts         MemberSyncOpts   `json:"adminSyncOpts"         description:"Settings for syncing system admin users"`
	SyncSchedule          string           `json:"syncSchedule"          description:"The sync job schedule in CRON format"`
}

// MakeLdapSettings returns a LDAPSettings response object for the given
// settings from the ldap backend.
func MakeLdapSettings(settings *ldapconfig.Settings) LDAPSettings {
	return LDAPSettings{
		RecoveryAdminUsername: settings.RecoveryAdminUsername,
		ServerURL:             settings.ServerURL,
		NoSimplePagination:    settings.NoSimplePagination,
		StartTLS:              settings.StartTLS,
		RootCerts:             settings.RootCerts,
		TLSSkipVerify:         settings.TLSSkipVerify,
		ReaderDN:              settings.ReaderDN,
		ReaderPassword:        settings.ReaderPassword,
		UserSearchConfigs:     makeUserSearchConfigs(settings.UserSearchConfigs),
		AdminSyncOpts:         makeAdminSyncOpts(settings.AdminSyncOpts),
		SyncSchedule:          settings.SyncSchedule,
	}
}

// UserSearchOpts is used in LDAPSettings.
type UserSearchOpts struct {
	BaseDN       string `json:"baseDN"       description:"The distinguished name of the element from which the LDAP server will search for users"`
	ScopeSubtree bool   `json:"scopeSubtree" description:"Whether to search for users in the entire subtree of the base DN or to only search one level under the base DN (if false)"`
	UsernameAttr string `json:"usernameAttr" description:"The name of the attribute of the LDAP user element which should be selected as the username"`
	FullNameAttr string `json:"fullNameAttr" description:"The name of the attribute of the LDAP user element which should be selected as the full name of the user"`
	Filter       string `json:"filter"       description:"The LDAP search filter used to select user elements, may be left blank"`
}

// makeUserSearchConfigs returns a slice of UserSearchOpts response objects
// from the given slice of options from the ldap backend.
func makeUserSearchConfigs(optsList []ldapconfig.UserSearchOpts) []UserSearchOpts {
	output := make([]UserSearchOpts, len(optsList))
	for i, opts := range optsList {
		output[i] = UserSearchOpts{
			BaseDN:       opts.BaseDN,
			ScopeSubtree: opts.ScopeSubtree,
			UsernameAttr: opts.UsernameAttr,
			FullNameAttr: opts.FullNameAttr,
			Filter:       opts.Filter,
		}
	}
	return output
}

// makeAdminSyncOpts returns a MemberSyncOpts response object for the given
// sync options from the ldap backend with sync enabled.
func makeAdminSyncOpts(opts ldapconfig.MemberSyncOpts) MemberSyncOpts {
	return MemberSyncOpts{
		EnableSync:         opts.EnableSync,
		SelectGroupMembers: opts.SelectGroupMembers,
		GroupDN:            opts.GroupDN,
		GroupMemberAttr:    opts.GroupMemberAttr,
		SearchBaseDN:       opts.SearchBaseDN,
		SearchScopeSubtree: opts.SearchScopeSubtree,
		SearchFilter:       opts.SearchFilter,
	}
}

// A Job object contains fields for a worker job.
type Job struct {
	ID          string    `json:"id"          description:"The ID of the job"`
	WorkerID    string    `json:"workerID"    description:"The ID of the worker which performed the job, unclaimed by a worker if empty"`
	Status      string    `json:"status"      description:"The current status of the job" enum:"waiting|running|done|canceled|errored"`
	ScheduledAt time.Time `json:"scheduledAt" description:"The time at which this job was scheduled"`
	LastUpdated time.Time `json:"lastUpdated" description:"The last time at which the status of this job was updated"`
	Action      string    `json:"action"      description:"The action this job performs"`
}

// MakeJob returns a Job response object for the given worker job.
func MakeJob(job *schema.Job) Job {
	return Job{
		ID:          job.ID,
		WorkerID:    job.WorkerID,
		Status:      job.Status,
		ScheduledAt: job.ScheduledAt,
		LastUpdated: job.LastUpdated,
		Action:      job.Action,
	}
}

// A Jobs object is used as a response for a list of jobs.
type Jobs struct {
	Jobs []Job `json:"jobs"`
}

// MakeJobs returns a Jobs response object for the given list of jobs.
func MakeJobs(jobs []schema.Job) Jobs {
	jobObjects := make([]Job, len(jobs))
	for i := range jobs {
		jobObjects[i] = MakeJob(&jobs[i])
	}

	return Jobs{
		Jobs: jobObjects,
	}
}

// A Worker represents a worker node in the jobs cluster.
type Worker struct {
	ID      string `json:"id"      description:"The ID of this worker"`
	Address string `json:"address" description:"Address at which an API server can contact the worker"`
}

// MakeWorker returns a Worker response object for the given worker.
func MakeWorker(worker *schema.Worker) Worker {
	return Worker{
		ID:      worker.ID,
		Address: worker.Address,
	}
}

// A Workers object is used as a response for a list of workers.
type Workers struct {
	Workers []Worker `json:"workers"`
}

// MakeWorkers returns a Workers response object for the given list of workers.
func MakeWorkers(workers []schema.Worker) Workers {
	workerObjects := make([]Worker, len(workers))
	for i := range workers {
		workerObjects[i] = MakeWorker(&workers[i])
	}

	return Workers{
		Workers: workerObjects,
	}
}

// JWK is used to represent a JSON Web Key.
type JWK struct {
	ID          string `json:"kid"           description:"Uniquely identifies a specific key"`
	KeyType     string `json:"kty"           description:"Specifies the type of the key as either RSA or Elliptic Curve" enum:"RSA|EC"`
	Modulus     string `json:"n,omitempty"   description:"Contains the modulus value for the RSA public key"`
	Exponent    string `json:"e,omitempty"   description:"Contains the exponent value for the RSA public key"`
	Curve       string `json:"crv,omitempty" description:"Identifies the cryptographic curve used with the key" enum:"P-256|P-384|P-521"`
	XCoordinate string `json:"x,omitempty"   description:"Contains the x coordinate for the Elliptic Curve point"`
	YCoordinate string `json:"y,omitempty"   description:"Contains the y coordinate for the Elliptic Curve point"`
}

// MakeJWK returns a JWK response object for the given key.
func MakeJWK(key *schema.JWK) JWK {
	return JWK{
		ID:          key.ID,
		KeyType:     key.KeyType,
		Modulus:     key.Modulus,
		Exponent:    key.Exponent,
		Curve:       key.Curve,
		XCoordinate: key.XCoordinate,
		YCoordinate: key.YCoordinate,
	}
}

// JWKSet represents a list of JWKs.
type JWKSet struct {
	Keys []JWK `json:"keys"`
}

// MakeJWKSet returns a JWKSet response object for the given list of keys.
func MakeJWKSet(keys []schema.JWK) JWKSet {
	keyObjects := make([]JWK, len(keys))
	for i := range keys {
		keyObjects[i] = MakeJWK(&keys[i])
	}

	return JWKSet{
		Keys: keyObjects,
	}
}

// Service respresents a service which relies on the auth service.
type Service struct {
	ID                 string   `json:"id"                 description:"uniqely identifies a service"`
	OwnerID            string   `json:"ownerID"            description:"ID of the account which owns and manages this service"`
	Name               string   `json:"name"               description:"simple name of this service; unique to the owning account"`
	Description        string   `json:"description"        description:"short description of the service"`
	URL                string   `json:"url"                description:"web address for the service"`
	Privileged         bool     `json:"privileged"         description:"whether this service has implicit authorization to all accounts"`
	RedirectURIs       []string `json:"redirectURIs"       description:"a list of URLs which serve as callbacks from Oauth authorization requests"`
	JWKsURIs           []string `json:"jwksURIs"            description:"a list of URLs for the service's JSON Web Key set, the set of public keys it used to sign authentication tokens"`
	ProviderIdentities []string `json:"providerIdentities" description:"A list of identifiers that this service will use as the 'aud' (Audience) for its authentication tokens (the tokens must have at least one 'aud' value that matches). Identity tokens issued by the provider will also use the first matching 'aud' (Audience) value in the 'iss' (Issuer) field"`
	CABundle           string   `json:"caBundle"           description:"PEM certificate bundle used to authenticate the JWKs URL endpoint"`
}

// MakeService returns a service response object for the given service.
func MakeService(service *schema.Service) Service {
	return Service{
		ID:                 service.ID,
		OwnerID:            service.OwnerID,
		Name:               service.Name,
		Description:        service.Description,
		URL:                service.URL,
		Privileged:         service.Privileged,
		RedirectURIs:       service.RedirectURIs,
		JWKsURIs:           service.JWKsURIs,
		ProviderIdentities: service.ProviderIdentities,
		CABundle:           service.CABundle,
	}
}

// A Services object is used as a response for a list of services.
type Services struct {
	Services []Service `json:"services"`
}

// MakeServices returns a Services response object for the given list of
// services.
func MakeServices(services []schema.Service) Services {
	serviceObjects := make([]Service, len(services))
	for i := range services {
		serviceObjects[i] = MakeService(&services[i])
	}

	return Services{
		Services: serviceObjects,
	}
}

// A LoginSession object is used as a response for a session login request.
type LoginSession struct {
	Account      Account `json:"account"      description:"the authenticated user account"`
	SessionToken string  `json:"sessionToken" description:"the session token created by the login"`
}

// AuthenticateRequest sets a session token authorization header on the given
// request.
func (s *LoginSession) AuthenticateRequest(r *http.Request) {
	r.Header.Set("Authorization", fmt.Sprintf("SessionToken %s", s.SessionToken))
}

// MakeLoginSession returns a LoginSession response object for the given user
// account and session token.
func MakeLoginSession(user *schema.Account, sessionToken string) LoginSession {
	return LoginSession{
		Account:      MakeAccount(user),
		SessionToken: sessionToken,
	}
}
