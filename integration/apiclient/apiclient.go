package apiclient

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy/adminserver"
	"github.com/docker/dhe-deploy/adminserver/api/common/forms"
	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/tmpforms"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/tmpresponses"

	distributionconfig "github.com/docker/distribution/configuration"
	enziclient "github.com/docker/orca/enzi/api/client"
	enziresponses "github.com/docker/orca/enzi/api/responses"
	"github.com/gorilla/websocket"
)

type APIClient interface {
	EnziSession() *enziclient.Session
	Login(username, password string) error
	Logout() error
	GetHost() string
	GetApiClientPort() uint16
	SetApiClientPort(uint16)
	GetApiClientUrlScheme() string
	SetApiClientUrlScheme(string)

	// Info
	Version() (string, error)

	// Accounts
	// CreateUser(username, password string) (*Account, error)
	// CreateUserWithLDAPLogin(username, ldapLogin, password string) (*Account, error)
	// CreateOrganization(name string) (*Account, error)
	// ListAccounts() ([]*Account, error)
	// ListUserAccounts() ([]*Account, error)
	// ListOrgAccounts() ([]*Account, error)
	GetAccount(name string) (*Account, error)
	// DeleteAccount(name string) error
	// ChangePassword(username, oldPassword, newPassword string) error
	// ActivateUser(username string) error
	// DeactivateUser(username string) error
	// ListOrganizations(username string) ([]*Account, error)

	// Teams
	// CreateManagedTeam(orgname, teamname, description string) (*Team, error)
	// CreateLDAPGroupSyncedTeam(orgname, teamname, description, groupDN, memberAttr string) (*Team, error)
	// ListTeams(orgname string) ([]*Team, error)
	// GetTeam(orgname, teamname string) (*Team, error)
	// UpdateTeam(orgname, teamname string, form TeamUpdateForm) (*Team, error)
	// DeleteTeam(orgname, teamname string) error
	// ListTeamMembers(orgname, teamname string) ([]*Account, error)
	// ListOrgMembers(orgname string) ([]*Account, error)
	// ListOrganizationMemberTeams(orgname, username string) ([]*Team, error)
	// CheckTeamMembership(orgname, teamname, username string) (isMember bool, err error)
	// CheckOrganizationMembership(orgname, username string) (isMember bool, err error)
	// ListUserOrganizations(username string) ([]*Account, error)
	// AddTeamMember(orgname, teamname, username string) error
	// DeleteTeamMember(orgname, teamname, username string) error
	// DeleteOrganizationMember(orgname, username string) error

	// Repository
	CreateRepository(namespace, reponame, shortDescription, longDescription, visibility string) (*responses.Repository, error)
	ListAllRepositories() ([]*responses.Repository, error)
	ListRepositories(namespace string) ([]*responses.Repository, error)
	// ListSharedRepositories(username string) ([]*responses.Repository, error)
	GetRepository(namespace, reponame string) (*responses.Repository, error)
	UpdateRepository(namespace, reponame string, form RepositoryUpdateForm) (*responses.Repository, error)
	DeleteRepository(namespace, reponame string) error

	// Tags
	GetRepositoryTags(namespace, reponame string) ([]responses.Tag, error)
	GetTagTrust(namespace, reponame, tag string) (*responses.Tag, error)
	DeleteTag(namespace, reponame, reference string) error

	// Manifests
	GetRepositoryManifests(namespace, reponame string) ([]responses.Manifest, error)
	DeleteManifest(namespace, reponame, reference string) error

	// Repository Access
	// user accesses removed for 2.0
	GetUserRepositoryAccess(username, namespace, reponame string) (*responses.RepoUserAccess, error)
	// ListRepositoryUserAccess(namespace, reponame string) (*Repository, []*UserAccess, error)
	SetRepositoryUserAccess(namespace, reponame, username, accessLevel string) (*RepositoryUserAccess, error)
	RevokeRepositoryUserAccess(namespace, reponame, username string) error

	// GetUserRepositoryNamespaceAccess(username, namespace string) (*responses.RepoNamespaceUserAccess, error)

	ListRepositoryTeamAccess(namespace, reponame string) (*responses.Repository, []responses.TeamAccess, error)
	ListTeamRepositoryAccess(orgname, teamname string) (*responses.Team, []responses.RepoAccess, error)
	SetRepositoryTeamAccess(namespace, reponame, teamname, accessLevel string) (*responses.RepoTeamAccess, error)
	RevokeRepositoryTeamAccess(namespace, reponame, teamname string) error

	ListRepositoryNamespaceTeamAccess(namespace string) (*responses.Namespace, []responses.TeamAccess, error)
	GetRepositoryNamespaceTeamAccess(namespace, teamname string) (*responses.NamespaceTeamAccess, error)
	SetRepositoryNamespaceTeamAccess(namespace, teamname, accessLevel string) (*responses.NamespaceTeamAccess, error)
	RevokeRepositoryNamespaceTeamAccess(namespace, teamname string) error

	// Search
	Autocomplete(opts SearchOptions) (*responses.Autocomplete, error)
	Reindex() error

	// HTTP
	GetHTTPSettings() (*responses.Settings, error)
	GetHTTPSettingsResponse() (*http.Response, error)
	SetHTTPSettings(settings *forms.Settings) error
	SetHTTPSettingsResponse(settings *forms.Settings) (*http.Response, error)

	// Registry
	GetRegistrySettings() (*distributionconfig.Configuration, error)
	SetRegistrySettings(*distributionconfig.Configuration) error

	GetCA() ([]byte, error)
	GetClusterStatus() (*responses.ClusterStatus, error)
	GetEvents() (*responses.Events, *url.Values, error)
	GetEventsWithParams(url.Values) (*responses.Events, *url.Values, error)
	EventsWebsocket() (*websocket.Conn, error)
	EventsWebsocketAuthd() (*websocket.Conn, error)

	// Auth
	// LDAPCheck(*adminserver.LDAPCheckSettings) (*map[string]interface{}, error)
	// GetAuthSettings() (*adminserver.AuthSettings, error)
	// GetAuthSettingsResponse() (*http.Response, error)
	// SetAuthSettings(settings *adminserver.AuthSettings) (*adminserver.AuthSettings, error)
	// SetAuthSettingsResponse(settings *adminserver.AuthSettings) (*http.Response, error)
	// ListSqlUsers() ([]*adminserver.FormUser, error)
	// ListSqlUsersResponse() (*http.Response, error)

	// License
	GetLicenseSettings() (*adminserver.LicenseSettings, error)
	SetLicenseSettings(licenseConfig *hubconfig.LicenseConfig) (*adminserver.LicenseSettings, error)
	SetLicenceSetttingsResponse(licenseConfig *hubconfig.LicenseConfig) (*http.Response, error)

	// Load balancer
	LoadBalancerStatus() (*NginxLoadBalancerStatus, error)
	LoadBalancerStatusResponse() (*http.Response, error)

	// Upgrade
	Upgrade(toVersion string) error

	// job framework
	GetJobStatus(id string) (string, error)
	RunJob(job tmpforms.JobSubmission) (*tmpresponses.Job, error)
	RunJobByAction(action string) (*tmpresponses.Job, error)

	// Set, get and delete schedules for garbage collection
	SetGarbageCollectionSchedule(schedule string) error
	GetGarbageCollectionSchedule() (string, error)
	DeleteGarbageCollectionSchedule() error

	Username() string
	// Password() string
}

type apiClient struct {
	loginSession  *enziresponses.LoginSession
	sessionSecret string
	csrfCookie    string

	enziHost  string
	ucpAsEnzi bool

	host               string
	apiClientPort      uint16
	apiClientUrlScheme string
	client             *http.Client
	retryAttempts      int
}

// ensure apiClient implements APIClient
var _ APIClient = (*apiClient)(nil)

func (c *apiClient) GetHost() string {
	return c.host
}

func (c *apiClient) GetApiClientPort() uint16 {
	return c.apiClientPort
}

func (c *apiClient) SetApiClientPort(newApiClientPort uint16) {
	c.apiClientPort = newApiClientPort
}

func (c *apiClient) GetApiClientUrlScheme() string {
	return c.apiClientUrlScheme
}

func (c *apiClient) SetApiClientUrlScheme(newApiClientUrlScheme string) {
	c.apiClientUrlScheme = newApiClientUrlScheme
}

func New(host string, retryAttempts int, httpClient *http.Client) APIClient {
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return RedirectAttemptedError
	}
	port := uint16(443)
	parts := strings.Split(host, ":")
	if len(parts) > 1 {
		host = parts[0]
		tmp, err := strconv.ParseInt(parts[1], 10, 16)
		if err != nil {
			logrus.Errorf("Failed to parse port: %s", err)
			port = 443
		} else {
			port = uint16(tmp)
		}
	}
	return &apiClient{
		host:               host,
		retryAttempts:      retryAttempts,
		client:             httpClient,
		apiClientUrlScheme: "https",
		apiClientPort:      port,
	}
}
