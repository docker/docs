package authz

import (
	"fmt"
	"strings"

	"github.com/docker/dhe-deploy/garant/authn"
	"github.com/docker/dhe-deploy/garant/authn/enzi"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/util"
	"github.com/docker/dhe-deploy/manager/schema"
	"github.com/docker/dhe-deploy/shared/dtrutil"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/distribution/context"
	garantauth "github.com/docker/garant/auth"
	"github.com/docker/garant/auth/common"
	"github.com/docker/orca/enzi/api/client/openid"
	enziresponses "github.com/docker/orca/enzi/api/responses"
)

var (
	AccessLevelScopeSets = map[string]common.ScopeSet{
		schema.AccessLevelAdmin:     common.NewScopeSet("*"),
		schema.AccessLevelReadWrite: common.NewScopeSet("pull", "push"),
		schema.AccessLevelReadOnly:  common.NewScopeSet("pull"),
	}
)

// Authorize should attempt to authorize the given account for the
// requested access to resources hosted by the specified service. The
// authorizer should return a subset of the requested access set which
// represents the access which the account has been granted. The authorizer
// should ignore access requests for resources that do not exist or for
// actions that do not exist in the context of the resource type. Any
// non-nil error returned will be interpreted as a server error.
func (a *authorizer) Authorize(ctx context.Context, acct garantauth.Account, service string, requestedAccess ...garantauth.Access) (grantedAccess []garantauth.Access, err error) {
	user, ok := acct.(*authn.User)
	if !ok || user == nil {
		// Must be an account authenticated with the managed auth backend.
		return nil, common.WithStackTrace(fmt.Errorf("invalid or nil account: %#v", acct))
	}

	requestedRepoScopeSets := common.GetValidRepositoryScopeSets(requestedAccess)
	repoAccessList := common.MakeRepoAccessList(a.filterAccessSetsForUser(requestedRepoScopeSets, user))
	if *user.Account.IsAdmin {
		// add in catalog access, if requested.
		catalogAccess := garantauth.Access{
			Action: "*",
			Resource: garantauth.Resource{
				Type: "registry",
				Name: "catalog",
			},
		}

		for _, access := range requestedAccess {
			if access == catalogAccess {
				repoAccessList = append(repoAccessList, access)
			}
		}
	}

	// add in search access, if requested.
	searchAccess := garantauth.Access{
		Action: "search",
		Resource: garantauth.Resource{
			Type: "registry",
			Name: "catalog",
		},
	}

	for _, access := range requestedAccess {
		if access == searchAccess {
			repoAccessList = append(repoAccessList, access)
		}
	}

	return repoAccessList, nil
}

func (a *authorizer) filterAccessSetsForUser(requestedRepoScopeSets common.RepositoryScopeSets, user *authn.User) common.RepositoryScopeSets {
	filteredAccessSets := make(common.RepositoryScopeSets, len(requestedRepoScopeSets))

	// Get each repository from the database. Note, this makes a db query per
	// each requested repository, however this method is typically only ever
	// called with a single requested repository.
	for repoName, requestedScopeSet := range requestedRepoScopeSets {
		tokens := strings.Split(repoName, "/")
		if len(tokens) < 2 {
			log.Errorf("repoName %q has no namespace", repoName)
			continue
		}
		namespaceName, name := strings.Join(tokens[:len(tokens)-1], "/"), tokens[len(tokens)-1]
		// ignore the dtr hostname segment if there is one
		if len(tokens) > 2 && tokens[0] == a.dtrHost {
			namespaceName = strings.Join(tokens[1:len(tokens)-1], "/")
		}

		ns, err := a.getNamespace(namespaceName, user)
		if err != nil {
			log.Errorf("unable to get namespace for repo %q", repoName)
			continue
		}

		repository, err := a.repoMgr.GetRepositoryByName(ns.ID, name)
		if err == schema.ErrNoSuchRepository {
			if *user.Account.IsAdmin {
				// admins get read-only access to unlisted repos
				// TODO remove this in 1.5?
				filteredAccessSets[repoName] = AccessLevelScopeSets[schema.AccessLevelReadOnly]
			}
			continue
		} else if err != nil {
			// Log the error and continue.
			log.Errorf("unable to get requested repository %q: %s", repoName, err)
			continue
		}

		accessLevel, err := a.RepositoryAccess(user, repository, ns)
		if err != nil {
			// Log the error and continue.
			log.Errorf("unable to get access level for repository %q for user %q: %s", repoName, user.Account.Name, err)
			continue
		}

		filteredAccessSets[repoName] = requestedScopeSet.Intersect(AccessLevelScopeSets[accessLevel])
	}

	return filteredAccessSets
}

func getOpenIDClient(settingsStore hubconfig.SettingsStore) (*openid.Client, error) {
	haConfig, err := settingsStore.HAConfig()
	if err != nil {
		return nil, err
	}

	enziConfig := util.GetEnziConfig(haConfig)
	client, err := dtrutil.HTTPClient(!enziConfig.VerifyCert, enziConfig.CA)
	if err != nil {
		return nil, err
	}

	return enzi.NewOpenIDClient(client, settingsStore)
}

func (a *authorizer) getNamespace(ns string, user *authn.User) (*enziresponses.Account, error) {
	if !user.IsAnonymous {
		return user.EnziSession.GetAccount(ns)
	}

	// NOTE: GetAccount is an authenticated endpoint. in order to
	// pull images as an anonymous user we need the account ID
	oc, err := getOpenIDClient(a.settingsStore)
	if err != nil {
		return nil, fmt.Errorf("unable to create a new openid client: %s", err)
	}

	// Get a token for the requested namespace. This will luckily have the namespace ID in it.
	// All is now good with the world, as we can use this ID from the external service
	// to look up data in our own database. Full steam ahead, el capitan.
	//
	// XXX: DO NOT use this token for ANYTHING else
	tokenResponse, err := oc.GetTokenWithRefreshToken(ns)
	if err != nil {
		return nil, fmt.Errorf("unable to create a request a token for the repository namespace %q: %s", ns, err)
	}
	return tokenResponse.Account, nil
}
