package ha_apiclient

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/docker/dhe-deploy/adminserver"
	"github.com/docker/dhe-deploy/adminserver/api/common/forms"
	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/integration/apiclient"
	"github.com/docker/dhe-deploy/jobs"
	distributionconfig "github.com/docker/distribution/configuration"
	enziclient "github.com/docker/orca/enzi/api/client"
)

// LoadBalancedAPIClient is an APIClient that emulates a load balancer
// across a set of apiclient.APIClients
type LoadBalancedAPIClient struct {
	loadBalancerPolicy string
	clients            []apiclient.APIClient

	// lastIdx is used to keep track of round robin load balancing
	lastIdx int
}

// Pick client returns one instance of an APIClient, determined by the policy
func (c *LoadBalancedAPIClient) pickClient() apiclient.APIClient {
	// Round Robin policy
	if c.loadBalancerPolicy == "roundRobin" {
		resClient := c.clients[c.lastIdx]
		c.lastIdx++
		if c.lastIdx == len(c.clients) {
			c.lastIdx = 0
		}
		return resClient
	}

	// Random policy
	return c.clients[rand.Intn(len(c.clients))]
}

func NewLoadBalancedAPIClient(subclients []apiclient.APIClient, policy string) apiclient.APIClient {
	rand.Seed(time.Now().UTC().UnixNano())
	if policy == "" {
		policy = "roundRobin"
	}
	return &LoadBalancedAPIClient{
		clients:            subclients,
		loadBalancerPolicy: policy,
		lastIdx:            0,
	}
}

// Method implementations - all methods use pickClient to route the request
// To one of the underlying APIClients

func (c *LoadBalancedAPIClient) EnziSession() *enziclient.Session {
	return c.pickClient().EnziSession()
}

func (c *LoadBalancedAPIClient) GetRegistrySettings() (*distributionconfig.Configuration, error) {
	return c.pickClient().GetRegistrySettings()
}

func (c *LoadBalancedAPIClient) SetRegistrySettings(cfg *distributionconfig.Configuration) error {
	return c.pickClient().SetRegistrySettings(cfg)
}

//func (c *LoadBalancedAPIClient) LoginWithSession(enziSession *enziclient.Session) error {
//	return c.pickClient().LoginWithSession(enziSession)
//}

func (c *LoadBalancedAPIClient) GetHost() string {
	return c.pickClient().GetHost()
}

func (c *LoadBalancedAPIClient) GetApiClientPort() uint16 {
	return c.pickClient().GetApiClientPort()
}

func (c *LoadBalancedAPIClient) SetApiClientPort(newApiClientPort uint16) {
	c.pickClient().SetApiClientPort(newApiClientPort)
}

func (c *LoadBalancedAPIClient) RevokeRepositoryUserAccess(namespace, reponame, username string) error {
	return c.pickClient().RevokeRepositoryUserAccess(namespace, reponame, username)
}

func (c *LoadBalancedAPIClient) SetRepositoryUserAccess(namespace, reponame, username, accesslevel string) (*apiclient.RepositoryUserAccess, error) {
	return c.pickClient().SetRepositoryUserAccess(namespace, reponame, username, accesslevel)
}

func (c *LoadBalancedAPIClient) GetApiClientUrlScheme() string {
	return c.pickClient().GetApiClientUrlScheme()
}

func (c *LoadBalancedAPIClient) SetApiClientUrlScheme(newApiClientUrlScheme string) {
	c.pickClient().SetApiClientUrlScheme(newApiClientUrlScheme)
}

func (c *LoadBalancedAPIClient) GetLicenseSettings() (*adminserver.LicenseSettings, error) {
	return c.pickClient().GetLicenseSettings()
}

func (c *LoadBalancedAPIClient) SetLicenseSettings(licenseConfig *hubconfig.LicenseConfig) (*adminserver.LicenseSettings, error) {
	return c.pickClient().SetLicenseSettings(licenseConfig)
}

// Really???? Three 't's in Setttings?
func (c *LoadBalancedAPIClient) SetLicenceSetttingsResponse(licenseConfig *hubconfig.LicenseConfig) (*http.Response, error) {
	return c.pickClient().SetLicenceSetttingsResponse(licenseConfig)
}

func (c *LoadBalancedAPIClient) Upgrade(toVersion string) error {
	return c.pickClient().Upgrade(toVersion)
}

//func (c *LoadBalancedAPIClient) CreateManagedTeam(orgname, teamname, description string) (*apiclient.Team, error) {
//	return c.pickClient().CreateManagedTeam(orgname, teamname, description)
//}

//func (c *LoadBalancedAPIClient) CreateLDAPGroupSyncedTeam(orgname, teamname, description, groupDN, memberAttr string) (*apiclient.Team, error) {
//	return c.pickClient().CreateLDAPGroupSyncedTeam(orgname, teamname, description, groupDN, memberAttr)
//}

//func (c *LoadBalancedAPIClient) ListTeams(orgname string) ([]*apiclient.Team, error) {
//	return c.pickClient().ListTeams(orgname)
//}

//func (c *LoadBalancedAPIClient) GetTeam(orgname, teamname string) (*apiclient.Team, error) {
//	return c.pickClient().GetTeam(orgname, teamname)
//}

//func (c *LoadBalancedAPIClient) UpdateTeam(orgname, teamname string, form apiclient.TeamUpdateForm) (*apiclient.Team, error) {
//	return c.pickClient().UpdateTeam(orgname, teamname, form)
//}

//func (c *LoadBalancedAPIClient) DeleteTeam(orgname, teamname string) error {
//	return c.pickClient().DeleteTeam(orgname, teamname)
//}

//func (c *LoadBalancedAPIClient) ListTeamMembers(orgname, teamname string) ([]*apiclient.Account, error) {
//	return c.pickClient().ListTeamMembers(orgname, teamname)
//}

//func (c *LoadBalancedAPIClient) ListOrganizationMemberTeams(orgname, username string) ([]*apiclient.Team, error) {
//	return c.pickClient().ListOrganizationMemberTeams(orgname, username)
//}

//func (c *LoadBalancedAPIClient) ListOrgMembers(orgname string) ([]*apiclient.Account, error) {
//	return c.pickClient().ListOrgMembers(orgname)
//}

//func (c *LoadBalancedAPIClient) ListUserOrganizations(username string) ([]*apiclient.Account, error) {
//	return c.pickClient().ListUserOrganizations(username)
//}

//func (c *LoadBalancedAPIClient) CheckTeamMembership(orgname, teamname, username string) (isMember bool, err error) {
//	return c.pickClient().CheckTeamMembership(orgname, teamname, username)
//}

//func (c *LoadBalancedAPIClient) CheckOrganizationMembership(orgname, username string) (isMember bool, err error) {
//	return c.pickClient().CheckOrganizationMembership(orgname, username)
//}

//func (c *LoadBalancedAPIClient) AddTeamMember(orgname, teamname, username string) error {
//	return c.pickClient().AddTeamMember(orgname, teamname, username)
//}

//func (c *LoadBalancedAPIClient) DeleteTeamMember(orgname, teamname, username string) error {
//	return c.pickClient().DeleteTeamMember(orgname, teamname, username)
//}

//func (c *LoadBalancedAPIClient) DeleteOrganizationMember(orgname, username string) error {
//	return c.pickClient().DeleteOrganizationMember(orgname, username)
//}

func (c *LoadBalancedAPIClient) Login(username, password string) error {
	return c.pickClient().Login(username, password)
}

func (c *LoadBalancedAPIClient) Logout() error {
	return c.pickClient().Logout()
}

//func (c *LoadBalancedAPIClient) CreateUser(username, password string) (*apiclient.Account, error) {
//	return c.pickClient().CreateUser(username, password)
//}
//
//func (c *LoadBalancedAPIClient) CreateUserWithLDAPLogin(username, ldapLogin, password string) (*apiclient.Account, error) {
//	return c.pickClient().CreateUserWithLDAPLogin(username, ldapLogin, password)
//}
//
//func (c *LoadBalancedAPIClient) CreateOrganization(name string) (*apiclient.Account, error) {
//	return c.pickClient().CreateOrganization(name)
//}
//
//func (c *LoadBalancedAPIClient) ListAccounts() ([]*apiclient.Account, error) {
//	return c.pickClient().ListAccounts()
//}
//
//func (c *LoadBalancedAPIClient) ListUserAccounts() ([]*apiclient.Account, error) {
//	return c.pickClient().ListUserAccounts()
//}
//
//func (c *LoadBalancedAPIClient) ListOrgAccounts() ([]*apiclient.Account, error) {
//	return c.pickClient().ListOrgAccounts()
//}
//
//func (c *LoadBalancedAPIClient) GetAccount(name string) (*apiclient.Account, error) {
//	return c.pickClient().GetAccount(name)
//}
//
//func (c *LoadBalancedAPIClient) DeleteAccount(name string) error {
//	return c.pickClient().DeleteAccount(name)
//}
//
//func (c *LoadBalancedAPIClient) ChangePassword(username, oldPassword, newPassword string) error {
//	return c.pickClient().ChangePassword(username, oldPassword, newPassword)
//}
//
//func (c *LoadBalancedAPIClient) ActivateUser(username string) error {
//	return c.pickClient().ActivateUser(username)
//}
//
//func (c *LoadBalancedAPIClient) DeactivateUser(username string) error {
//	return c.pickClient().DeactivateUser(username)
//}
//
//func (c *LoadBalancedAPIClient) ListOrganizations(username string) ([]*apiclient.Account, error) {
//	return c.pickClient().ListOrganizations(username)
//}
//
func (c *LoadBalancedAPIClient) Username() string {
	return c.pickClient().Username()
}

//func (c *LoadBalancedAPIClient) Password() string {
//	return c.pickClient().Password()
//}

//func (c *LoadBalancedAPIClient) LDAPCheck(params *adminserver.LDAPCheckSettings) (*map[string]interface{}, error) {
//	return c.pickClient().LDAPCheck(params)
//}
//
//func (c *LoadBalancedAPIClient) GetAuthSettings() (*adminserver.AuthSettings, error) {
//	return c.pickClient().GetAuthSettings()
//}
//
//func (c *LoadBalancedAPIClient) GetAuthSettingsResponse() (*http.Response, error) {
//	return c.pickClient().GetAuthSettingsResponse()
//}
//
//func (c *LoadBalancedAPIClient) SetAuthSettings(settings *adminserver.AuthSettings) (*adminserver.AuthSettings, error) {
//	return c.pickClient().SetAuthSettings(settings)
//}
//
//func (c *LoadBalancedAPIClient) SetAuthSettingsResponse(settings *adminserver.AuthSettings) (*http.Response, error) {
//	return c.pickClient().SetAuthSettingsResponse(settings)
//}
//
//func (c *LoadBalancedAPIClient) ListSqlUsers() ([]*adminserver.FormUser, error) {
//	return c.pickClient().ListSqlUsers()
//}
//
//func (c *LoadBalancedAPIClient) ListSqlUsersResponse() (*http.Response, error) {
//	return c.pickClient().ListSqlUsersResponse()
//}
//

//func (c *LoadBalancedAPIClient) GetAuthSettings() (*adminserver.AuthSettings, error) {
//	return c.pickClient().GetAuthSettings()
//}

//func (c *LoadBalancedAPIClient) GetAuthSettingsResponse() (*http.Response, error) {
//	return c.pickClient().GetAuthSettingsResponse()
//}

//func (c *LoadBalancedAPIClient) SetAuthSettings(settings *adminserver.AuthSettings) (*adminserver.AuthSettings, error) {
//	return c.pickClient().SetAuthSettings(settings)
//}

//func (c *LoadBalancedAPIClient) SetAuthSettingsResponse(settings *adminserver.AuthSettings) (*http.Response, error) {
//	return c.pickClient().SetAuthSettingsResponse(settings)
//}

//func (c *LoadBalancedAPIClient) ListSqlUsers() ([]*adminserver.FormUser, error) {
//	return c.pickClient().ListSqlUsers()
//}

//func (c *LoadBalancedAPIClient) ListSqlUsersResponse() (*http.Response, error) {
//	return c.pickClient().ListSqlUsersResponse()
//}

func (c *LoadBalancedAPIClient) SetGarbageCollectionSchedule(schedule string) error {
	return c.pickClient().SetGarbageCollectionSchedule(schedule)
}

func (c *LoadBalancedAPIClient) DeleteGarbageCollectionSchedule() error {
	return c.pickClient().DeleteGarbageCollectionSchedule()
}

func (c *LoadBalancedAPIClient) GetGarbageCollectionSchedule() (string, error) {
	return c.pickClient().GetGarbageCollectionSchedule()
}

func (c *LoadBalancedAPIClient) GetHTTPSettings() (*responses.Settings, error) {
	return c.pickClient().GetHTTPSettings()
}

func (c *LoadBalancedAPIClient) GetHTTPSettingsResponse() (*http.Response, error) {
	return c.pickClient().GetHTTPSettingsResponse()
}

func (c *LoadBalancedAPIClient) SetHTTPSettings(settings *forms.Settings) error {
	return c.pickClient().SetHTTPSettings(settings)
}

func (c *LoadBalancedAPIClient) SetHTTPSettingsResponse(settings *forms.Settings) (*http.Response, error) {
	return c.pickClient().SetHTTPSettingsResponse(settings)
}

func (c *LoadBalancedAPIClient) GetJobStatus(job string) (*jobs.JobStatus, error) {
	return c.pickClient().GetJobStatus(job)
}

func (c *LoadBalancedAPIClient) RunJob(job string) error {
	return c.pickClient().RunJob(job)
}

func (c *LoadBalancedAPIClient) LoadBalancerStatus() (*apiclient.NginxLoadBalancerStatus, error) {
	return c.pickClient().LoadBalancerStatus()
}

func (c *LoadBalancedAPIClient) LoadBalancerStatusResponse() (*http.Response, error) {
	return c.pickClient().LoadBalancerStatusResponse()
}

func (c *LoadBalancedAPIClient) CreateRepository(namespace, reponame, shortDescription, longDescription, visibility string) (*responses.Repository, error) {
	return c.pickClient().CreateRepository(namespace, reponame, shortDescription, longDescription, visibility)
}

//func (c *LoadBalancedAPIClient) ListSharedRepositories(username string) ([]*apiclient.Repository, error) {
//	return c.pickClient().ListSharedRepositories(username)
//}

func (c *LoadBalancedAPIClient) ListAllRepositories() ([]*responses.Repository, error) {
	return c.pickClient().ListAllRepositories()
}

func (c *LoadBalancedAPIClient) ListRepositories(namespace string) ([]*responses.Repository, error) {
	return c.pickClient().ListRepositories(namespace)
}

func (c *LoadBalancedAPIClient) GetRepository(namespace, reponame string) (*responses.Repository, error) {
	return c.pickClient().GetRepository(namespace, reponame)
}

func (c *LoadBalancedAPIClient) UpdateRepository(namespace, reponame string, form apiclient.RepositoryUpdateForm) (*responses.Repository, error) {
	return c.pickClient().UpdateRepository(namespace, reponame, form)
}

func (c *LoadBalancedAPIClient) DeleteRepository(namespace, reponame string) error {
	return c.pickClient().DeleteRepository(namespace, reponame)
}

func (c *LoadBalancedAPIClient) ListRepositoryNamespaceTeamAccess(namespace string) (*responses.Namespace, []responses.TeamAccess, error) {
	return c.pickClient().ListRepositoryNamespaceTeamAccess(namespace)
}

func (c *LoadBalancedAPIClient) GetRepositoryNamespaceTeamAccess(namespace, teamname string) (*responses.NamespaceTeamAccess, error) {
	return c.pickClient().GetRepositoryNamespaceTeamAccess(namespace, teamname)
}

func (c *LoadBalancedAPIClient) SetRepositoryNamespaceTeamAccess(namespace, teamname, accessLevel string) (*responses.NamespaceTeamAccess, error) {
	return c.pickClient().SetRepositoryNamespaceTeamAccess(namespace, teamname, accessLevel)
}

func (c *LoadBalancedAPIClient) RevokeRepositoryNamespaceTeamAccess(namespace, teamname string) error {
	return c.pickClient().RevokeRepositoryNamespaceTeamAccess(namespace, teamname)
}

func (c *LoadBalancedAPIClient) ListTeamRepositoryAccess(orgname, teamname string) (*responses.Team, []responses.RepoAccess, error) {
	return c.pickClient().ListTeamRepositoryAccess(orgname, teamname)
}

func (c *LoadBalancedAPIClient) ListRepositoryTeamAccess(namespace, reponame string) (*responses.Repository, []responses.TeamAccess, error) {
	return c.pickClient().ListRepositoryTeamAccess(namespace, reponame)
}

func (c *LoadBalancedAPIClient) SetRepositoryTeamAccess(namespace, reponame, teamname, accessLevel string) (*responses.RepoTeamAccess, error) {
	return c.pickClient().SetRepositoryTeamAccess(namespace, reponame, teamname, accessLevel)
}

func (c *LoadBalancedAPIClient) RevokeRepositoryTeamAccess(namespace, reponame, teamname string) error {
	return c.pickClient().RevokeRepositoryTeamAccess(namespace, reponame, teamname)
}

func (c *LoadBalancedAPIClient) GetUserRepositoryAccess(username, namespace, reponame string) (*responses.RepoUserAccess, error) {
	return c.pickClient().GetUserRepositoryAccess(username, namespace, reponame)
}

//func (c *LoadBalancedAPIClient) GetUserRepositoryNamespaceAccess(username, namespace string) (*responses.RepoNamespaceUserAccess, error) {
//	return c.pickClient().GetUserRepositoryNamespaceAccess(username, namespace)
//}

//func (c *LoadBalancedAPIClient) ListRepositoryUserAccess(namespace, reponame string) (*apiclient.Repository, []*apiclient.UserAccess, error) {
//	return c.pickClient().ListRepositoryUserAccess(namespace, reponame)
//}

//func (c *LoadBalancedAPIClient) SetRepositoryUserAccess(namespace, reponame, username, accessLevel string) (*apiclient.RepositoryUserAccess, error) {
//	return c.pickClient().SetRepositoryUserAccess(namespace, reponame, username, accessLevel)
//}

//func (c *LoadBalancedAPIClient) RevokeRepositoryUserAccess(namespace, reponame, username string) error {
//	return c.pickClient().RevokeRepositoryUserAccess(namespace, reponame, username)
//}

//func (c *LoadBalancedAPIClient) Search(opts apiclient.SearchOptions) (*responses.Search, error) {
//	return c.pickClient().Search(opts)
//}

func (c *LoadBalancedAPIClient) Autocomplete(opts apiclient.SearchOptions) (*responses.Autocomplete, error) {
	return c.pickClient().Autocomplete(opts)
}

func (c *LoadBalancedAPIClient) Reindex() error {
	return c.pickClient().Reindex()
}

func (c *LoadBalancedAPIClient) GetRepositoryTags(namespace, reponame string) (*responses.ListRepositoryTags, error) {
	return c.pickClient().GetRepositoryTags(namespace, reponame)
}

func (c *LoadBalancedAPIClient) GetTagTrust(namespace, reponame, tag string) (*responses.TagWithSignature, error) {
	return c.pickClient().GetTagTrust(namespace, reponame, tag)
}

func (c *LoadBalancedAPIClient) DeleteManifestOrTag(namespace, reponame, reference string) error {
	return c.pickClient().DeleteManifestOrTag(namespace, reponame, reference)
}

func (c *LoadBalancedAPIClient) Version() (string, error) {
	return c.pickClient().Version()
}
