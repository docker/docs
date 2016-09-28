package authz

import (
	"errors"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/garant/authn"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/etcd"
	"github.com/docker/dhe-deploy/hubconfig/sanitizers"
	"github.com/docker/dhe-deploy/hubconfig/settingsstore"
	configutil "github.com/docker/dhe-deploy/hubconfig/util"
	"github.com/docker/dhe-deploy/manager/schema"
	"github.com/docker/dhe-deploy/shared/containers"
	"github.com/docker/dhe-deploy/shared/dtrutil"

	"github.com/docker/garant/auth"
	"github.com/docker/garant/config"
	enziresponses "github.com/docker/orca/enzi/api/responses"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

// Authorizer defines an interface with all of the methods we need for
// being a garant authorizer backend, authenticationg request users in DTR,
// and authorizing access to repositories in DTR.
type Authorizer interface {
	auth.TokenAuthorizer
	authn.Authenticator

	CheckAdminOrOrgOwner(user *authn.User, orgID string) error
	CheckAdminOrOrgMember(user *authn.User, orgID string) error
	NamespaceAccess(user *authn.User, ns *enziresponses.Account) (accessLevel string, err error)
	RepositoryAccess(user *authn.User, repo *schema.Repository, ns *enziresponses.Account) (accessLevel string, err error)
	AllVisibleRepositoriesInNamespace(user *authn.User, ns *enziresponses.Account) (visibleRepositories []*schema.Repository, err error)
	VisibleRepositoriesInNamespace(user *authn.User, ns *enziresponses.Account, startPK string, limit uint) (visibleRepositories []*schema.Repository, nextPK string, err error)
	FilterVisibleReposInNamespace(user *authn.User, namespace *enziresponses.Account, repos []*schema.Repository, isNamespaceAdmin bool) ([]*schema.Repository, error)
	SharedRepositoriesForUser(user *authn.User, startPK string, limit uint) (sharedRepositories []*schema.Repository, nextPK string, err error)
}

// ErrUninitialized indicates that auth has not been configured.
var ErrUninitialized = errors.New("uninitialized")

// Authorizer is capable of authenticating a user and determining access
// levels for that user.
type authorizer struct {
	enziHost string
	dtrHost  string

	settingsStore  hubconfig.SettingsStore
	repoMgr        *schema.RepositoryManager
	accessMgr      *schema.RepositoryAccessManager
	clientTokenMgr *schema.ClientTokenManager
}

// Assert that *authorizer implements the Authorizer interface.
var _ Authorizer = (*authorizer)(nil)

// NewAuthorizer creates a new authorizer using the given database connection.
func NewAuthorizer(session *rethink.Session, settingsStore hubconfig.SettingsStore) Authorizer {
	haConfig, err := settingsStore.HAConfig()
	if err != nil {
		logrus.WithField("error", err).Fatal("Failed to get HA config")
	}

	userHubConfig, err := settingsStore.UserHubConfig()
	if err != nil {
		logrus.WithField("error", err).Fatal("Failed to get UserHub config")
	}

	enziConfig := configutil.GetEnziConfig(haConfig)
	return &authorizer{
		enziHost: enziConfig.Host,
		dtrHost:  userHubConfig.DTRHost,

		settingsStore:  settingsStore,
		repoMgr:        schema.NewRepositoryManager(session),
		accessMgr:      schema.NewRepositoryAccessManager(session),
		clientTokenMgr: schema.NewClientTokenManager(session),
	}
}

func setupGarantAuthorizer(_ config.Parameters) (auth.Authorizer, error) {
	replicaID := os.Getenv(deploy.ReplicaIDEnvVar)
	session, err := dtrutil.GetRethinkSession(replicaID)
	if err != nil {
		return nil, err
	}

	kvStore, err := etcd.NewKeyValueStore(containers.EtcdUrls(), deploy.EtcdPath)
	if err != nil {
		logrus.WithField("error", err).Fatal("Failed to initialize key value storage")
	}

	return NewAuthorizer(session, sanitizers.Wrap(settingsstore.New(kvStore))), nil
}

func init() {
	auth.Register("dtr", auth.InitFunc(setupGarantAuthorizer))
}
