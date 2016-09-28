package mock_test

import (
	"bytes"
	"crypto"
	"io"
	"io/ioutil"

	"net/http"
	"net/url"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/engine-api/types/swarm"
	kvstore "github.com/docker/libkv/store"

	"github.com/docker/orca"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/controller/manager"
	"github.com/docker/orca/dockerhub"
	orcatypes "github.com/docker/orca/types"
)

type MockManager struct{}

func (m MockManager) Container(id string) (types.ContainerJSON, error) {
	// switch labels if needed
	labels := map[string]string{}
	switch id {
	case "none":
		labels[orca.UCPAccessLabel] = "none"
	case "view":
		labels[orca.UCPAccessLabel] = "view"
	case "restricted":
		labels[orca.UCPAccessLabel] = "restricted"
	case "full":
		labels[orca.UCPAccessLabel] = "full"
	case TestContainerAccessLabel:
		labels[orca.UCPAccessLabel] = TestContainerAccessLabel
	}

	// owner label
	labels[orca.UCPOwnerLabel] = TestContainerAccessOwner

	return getTestContainerInfo(TestContainerId, TestContainerName, TestContainerImage, labels), nil
}

func (m MockManager) Service(id string) (swarm.Service, error) {
	return swarm.Service{}, nil
}

func (m MockManager) Network(id string) (types.NetworkResource, error) {
	return types.NetworkResource{}, nil
}

func (m MockManager) ListUserNetworks(ctx *auth.Context, filters filters.Args) ([]types.NetworkResource, error) {
	return []types.NetworkResource{}, nil
}

func (m MockManager) ListServices(filters filters.Args) ([]swarm.Service, error) {
	return []swarm.Service{}, nil
}

func (m MockManager) ListUserServices(ctx *auth.Context, filters filters.Args) ([]swarm.Service, error) {
	return []swarm.Service{}, nil
}

func (m MockManager) ListContainers(all bool, size bool, filters string) ([]types.Container, error) {
	return TestContainers, nil
}

func (m MockManager) ListUserContainers(ctx *auth.Context, all bool, size bool, filters string) ([]types.Container, error) {
	return TestContainers, nil
}
func (m MockManager) ContainerLogs(id string, opts types.ContainerLogsOptions) (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBuffer([]byte("test log message"))), nil
}

func (m MockManager) DockerClient() *client.Client {
	client, _ := client.NewClient("unix:///var/run/docker.sock", "", nil, nil)
	return client
}
func (m MockManager) DockerClientTransport() *http.Transport {
	return &http.Transport{}
}

func (m MockManager) ProxyClient() *client.Client {
	return nil
}
func (m MockManager) ProxyClientTransport() *http.Transport {
	return nil
}

func (m MockManager) SwarmClassicURL() *url.URL {
	return nil
}

func (m MockManager) EngineProxyURL() *url.URL {
	return nil
}

func (M MockManager) EnableAdminUCPScheduling() bool {
	return true
}

func (M MockManager) EnableUserUCPScheduling() bool {
	return true
}

func (m MockManager) GenerateClientBundle(ctx *auth.Context, host, label string) ([]byte, error) {
	return nil, nil
}

func (m MockManager) ID() string {
	return "mock-manager"
}

func (m MockManager) Datastore() kvstore.Store {
	return nil
}

func (m MockManager) SaveEvent(event *orca.Event) error {
	return nil
}

func (m MockManager) Accounts(ctx *auth.Context) ([]*auth.Account, error) {
	return []*auth.Account{
		TestAccount,
	}, nil
}

func (m MockManager) Account(ctx *auth.Context, username string) (*auth.Account, error) {
	if username == "admin" {
		return &auth.Account{
			Username: "admin",
			Admin:    true,
		}, nil
	} else if username == TestAccount.Username {
		return TestAccount, nil
	}
	return nil, auth.ErrAccountDoesNotExist
}

func (m MockManager) AccountsForTeam(ctx *auth.Context, teamID string) ([]*auth.Account, error) {
	return nil, nil
}

func (m MockManager) SaveAccount(ctx *auth.Context, updateAccount *auth.Account) (string, error) {
	return "", nil
}

func (m MockManager) DeleteAccount(ctx *auth.Context, account *auth.Account) error {
	return nil
}

func (m MockManager) AuthenticateUsernamePassword(username, password, remoteAddr string) (*auth.Context, error) {
	return nil, nil
}

func (m MockManager) AuthenticatePublicKey(key crypto.PublicKey, remoteAddr string) (*auth.Context, error) {
	return nil, nil
}

func (m MockManager) AuthenticateSessionToken(token string) (*auth.Context, error) {
	return nil, manager.ErrInvalidAuthToken
}

func (m MockManager) ChangePassword(ctx *auth.Context, username, oldPassword, newPassword string) error {
	return nil
}

func (m MockManager) WebhookKeys() ([]*dockerhub.WebhookKey, error) {
	return []*dockerhub.WebhookKey{
		TestWebhookKey,
	}, nil
}

func (m MockManager) NewWebhookKey(image string) (*dockerhub.WebhookKey, error) {
	return nil, nil
}

func (m MockManager) WebhookKey(key string) (*dockerhub.WebhookKey, error) {
	return nil, nil
}

func (m MockManager) SaveWebhookKey(key *dockerhub.WebhookKey) error {
	return nil
}

func (m MockManager) DeleteWebhookKey(id string) error {
	return nil
}

func (m MockManager) AddRegistry(registry orca.Registry) error {
	return nil
}

func (m MockManager) Registries() ([]orca.Registry, error) {
	return []orca.Registry{
		TestRegistry,
	}, nil
}

func (m MockManager) Registry(name string) (orca.Registry, error) {
	return TestRegistry, nil
}

func (m MockManager) RemoveRegistry(registry orca.Registry) error {
	return nil
}

func (m MockManager) Nodes() ([]*orca.Node, error) {
	return []*orca.Node{
		TestNode,
	}, nil
}

func (m MockManager) AuthorizeNodeRequest(r *orca.NodeRequest) (*orca.NodeConfiguration, error) {
	return nil, nil
}

func (m MockManager) DeleteRepository(name string) error {
	return nil
}

func (m MockManager) Node(name string) (*orca.Node, error) {
	return TestNode, nil
}

func (m MockManager) ListNodes() ([]swarm.Node, error) {
	return []swarm.Node{}, nil
}

func (m MockManager) InspectNode(nodeID string) (swarm.Node, error) {
	return swarm.Node{}, nil
}

func (m MockManager) PromoteNode(name string, replicate bool) error {
	return nil
}

func (m MockManager) SetControllerServerCerts(name string, req orcatypes.ServerCertRequest) error {
	return nil
}

func (m MockManager) CreateConsoleSession(c *orca.ConsoleSession) error {
	return nil
}

func (m MockManager) RemoveConsoleSession(c *orca.ConsoleSession) error {
	return nil
}

func (m MockManager) ConsoleSession(token string) (*orca.ConsoleSession, error) {
	return TestConsoleSession, nil
}

func (m MockManager) ValidateConsoleSessionToken(containerId, token string) bool {
	return true
}

func (m MockManager) CreateContainerLogsToken(c *orca.ContainerLogsToken) error {
	return nil
}

func (m MockManager) RemoveContainerLogsToken(c *orca.ContainerLogsToken) error {
	return nil
}

func (m MockManager) ContainerLogsToken(token string) (*orca.ContainerLogsToken, error) {
	return TestContainerLogsToken, nil
}

func (m MockManager) ValidateContainerLogsToken(containerId, token string) bool {
	return true
}

func (m MockManager) GetTrackingDisabled() bool {
	return true
}

func (m MockManager) GetUsageInfoDisabled() bool {
	return true
}

func (m MockManager) AnonymizeTracking() bool {
	return true
}

func (m MockManager) TrackClientInfo(req *http.Request) {
	return
}

func (m MockManager) GetAuthenticator() auth.Authenticator {
	return nil
}

func (m MockManager) ScaleContainer(id string, numInstances int) manager.ScaleResult {
	return manager.ScaleResult{Scaled: []string{"9c3c7dd2199a95cce29950b612ecf918ae278a42e53e10f6cccb752b6fbcd8b3"}, Errors: []string{"500 Internal Server Error: no resources available to schedule container"}}
}

func (m MockManager) GetDefaultCatalog() ([]orca.CatalogItem, error) {
	return nil, nil
}

func (m MockManager) SearchCatalog(query string) ([]orca.CatalogItem, error) {
	return nil, nil
}

func (m MockManager) Applications(ctx *auth.Context) ([]*orca.Application, error) {
	return nil, nil
}

func (m MockManager) CreateApplication(application *orca.Application) error {
	return nil
}

func (m MockManager) Application(ctx *auth.Context, name string) (*orca.Application, error) {
	return nil, nil
}

func (m MockManager) RemoveApplication(application *orca.Application) error {
	return nil
}

func (m MockManager) DeployApplication(application *orca.Application) error {
	return nil
}

func (m MockManager) RestartApplication(application *orca.Application) error {
	return nil
}

func (m MockManager) StopApplication(application *orca.Application) error {
	return nil
}

func (m MockManager) ApplicationContainers(application *orca.Application) ([]types.Container, error) {
	return nil, nil
}

func (m MockManager) AddTrustedCert(caCert string) {
}

func (m MockManager) SupportDump(w io.Writer) error {
	return nil
}

func (m MockManager) Logout(tokenStr string) error {
	return nil
}

func (m MockManager) GetSelfStatus() error {
	return nil
}

func (m MockManager) GetStatus(node *manager.ManagerNode) string {
	return ""
}

func (m MockManager) GetManagers() []*manager.ManagerNode {
	return nil
}

func (m MockManager) ListConfigSubsystems() []string {
	return []string{}
}

func (m MockManager) GetSubsystemConfig(subsystemName string) (string, error) {
	return "", nil
}

func (m MockManager) UserConfigUpdate(subsystemName, jsonConfig string) error {
	return nil
}
func (m MockManager) GetLicense() manager.LicenseSubsystemConfig {
	return manager.LicenseSubsystemConfig{}
}
func (m MockManager) GetLicenseKeyID() string {
	return ""
}
func (m MockManager) GetLicenseTier() string {
	return ""
}
func (m MockManager) GetBanner(ctx *auth.Context) []manager.Banner {
	return nil
}
func (m MockManager) AccessLists() ([]*auth.AccessList, error) {
	return nil, nil
}
func (m MockManager) AccessListsForTeam(teamId string) ([]*auth.AccessList, error) {
	return nil, nil
}
func (m MockManager) AccessListsForAccount(ctx *auth.Context, username string) ([]*auth.AccessList, error) {
	return nil, nil
}
func (m MockManager) SaveAccessList(list *auth.AccessList) (string, error) {
	return "", nil
}
func (m MockManager) AccessList(teamId, id string) (*auth.AccessList, error) {
	return nil, nil
}
func (m MockManager) RemoveAccessList(teamId, id string) error {
	return nil
}
func (m MockManager) GetAccess(ctx *auth.Context) (map[string]auth.Role, error) {
	// TODO: switch based upon username
	return map[string]auth.Role{
		"none":       auth.None,
		"view":       auth.View,
		"restricted": auth.RestrictedControl,
		"full":       auth.FullControl,
		TestContainerAccessLabel: auth.FullControl,
	}, nil
}
func (m MockManager) ListImages(all bool, filter filters.Args) ([]types.Image, error) {
	return nil, nil
}
func (m MockManager) ListUserImages(ctx *auth.Context, all bool, filter filters.Args) ([]types.Image, error) {
	return nil, nil
}
func (m MockManager) Team(ctx *auth.Context, id string) (*auth.Team, error) {
	return TestTeam, nil
}
func (m MockManager) Teams(ctx *auth.Context) ([]*auth.Team, error) {
	return []*auth.Team{
		TestTeam,
	}, nil
}
func (m MockManager) SaveTeam(ctx *auth.Context, team *auth.Team) (string, error) {
	return "", nil
}
func (m MockManager) DeleteTeam(ctx *auth.Context, team *auth.Team) error {
	return nil
}
func (m MockManager) AddMemberToTeam(ctx *auth.Context, teamID, username string) error {
	return nil
}
func (m MockManager) RemoveMemberFromTeam(ctx *auth.Context, teamID, username string) error {
	return nil
}
func (m MockManager) AuthSyncMessages(ctx *auth.Context) string {
	return ""
}
func (m MockManager) AuthSync(ctx *auth.Context) string {
	return ""
}
func (m MockManager) RequireContentTrustForDTR() *bool {
	ret := false
	return &ret
}
func (m MockManager) RequireContentTrustForHub() *bool {
	ret := false
	return &ret
}
