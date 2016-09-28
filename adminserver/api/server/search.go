package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/adminserver/api/common/forms"
	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
	"github.com/docker/dhe-deploy/garant/authn"
	"github.com/docker/dhe-deploy/manager/schema"
	"github.com/docker/dhe-deploy/shared/dtrutil"

	"github.com/docker/distribution/context"
	enzierrors "github.com/docker/orca/enzi/api/errors"
	enziresponses "github.com/docker/orca/enzi/api/responses"
	"github.com/emicklei/go-restful"
)

const DockerSearchLimit = 25

func (a *APIServer) handleDockerSearch(ctx context.Context, r *restful.Request) responses.APIResponse {
	searchOpts := forms.SearchOptions{
		IncludeRepositories: true,
		Query:               r.QueryParameter("q"),
		Limit:               DockerSearchLimit,
	}

	user, _ := a.getAuthenticatedUser(r, false)
	_, repositoryResults, errResp := a.handleAutocompleteWithUserAndOptions(ctx, user, searchOpts)
	if errResp != nil {
		return errResp
	}
	return responses.JSONResponse(http.StatusOK, nil, nil, responses.MakeDockerSearch(repositoryResults))
}

func (a *APIServer) handleAutocomplete(ctx context.Context, r *restful.Request) responses.APIResponse {
	user, _ := a.getAuthenticatedUser(r, false)
	options := parseSearchOptions(r)
	if options.Limit == 0 {
		options.Limit = 300 // arbitrary default limit
	}

	accountResults, repositoryResults, errResp := a.handleAutocompleteWithUserAndOptions(ctx, user, options)
	if errResp != nil {
		return errResp
	}
	return responses.JSONResponse(http.StatusOK, nil, nil, responses.MakeAutocomplete(accountResults, repositoryResults))
}

func (a *APIServer) handleAutocompleteWithUserAndOptions(ctx context.Context, user *authn.User, options forms.SearchOptions) ([]enziresponses.Account, []responses.Repository, responses.APIResponse) {
	var acctPrefix, repoPrefix string
	var accountResults []enziresponses.Account
	var accountsToListRepos []enziresponses.Account
	var failIfAccountNotExists bool
	var errResp responses.APIResponse
	if options.Namespace != "" {
		// if you specify the namespace then the entire query is for the
		// repository and we don't need to search accounts.
		acctPrefix = options.Namespace
		repoPrefix = options.Query
		failIfAccountNotExists = true
	} else {
		queryParts := strings.SplitN(options.Query, "/", 2)

		listReposForAccResults := true
		if len(queryParts) == 2 {
			if queryParts[0] == "" {
				listReposForAccResults = false
				repoPrefix = queryParts[1]
			} else if queryParts[1] == "" {
				acctPrefix = queryParts[0]
			} else {
				acctPrefix = queryParts[0]
				repoPrefix = queryParts[1]
			}
		} else {
			// there is 1 part
			acctPrefix = queryParts[0]
			repoPrefix = queryParts[0]
		}

		accountResults, errResp = a.autocompleteAccounts(ctx, user, acctPrefix, options.Limit)
		if errResp != nil {
			return nil, nil, errResp
		}
		if len(accountResults) > 0 && accountResults[0].Name == acctPrefix {
			// the acctPrefix is a real account, so make it the namespace
			// restriction for the repo autocomplete
			listReposForAccResults = false
			failIfAccountNotExists = true
		}
		if listReposForAccResults {
			accountsToListRepos = accountResults
		}
		if !options.IncludeAccounts {
			accountResults = nil
		}
	}

	var repositoryResults []responses.Repository
	if options.IncludeRepositories {
		repositoryResults, errResp = a.autocompleteRepositories(ctx, user, acctPrefix, repoPrefix, failIfAccountNotExists, accountsToListRepos, options.Limit)
		if errResp != nil {
			return nil, nil, errResp
		}
	}

	return accountResults, repositoryResults, nil
}

func (a *APIServer) autocompleteAccounts(ctx context.Context, user *authn.User, acctPrefix string, limit uint) ([]enziresponses.Account, responses.APIResponse) {
	if user == nil {
		// TODO after openID integration we should use that?
		return nil, nil
	}
	accounts, _, err := user.EnziSession.ListAccounts("", acctPrefix, limit)
	if err != nil {
		return nil, responses.APIError(errors.InternalError(ctx, err))
	}
	for i, account := range accounts.Accounts {
		if !strings.HasPrefix(account.Name, acctPrefix) {
			return accounts.Accounts[:i], nil
		}
	}
	return accounts.Accounts, nil
}

func (a *APIServer) autocompleteRepositories(ctx context.Context, user *authn.User, acctPrefix, repoPrefix string, failIfAccountNotExists bool, accountResults []enziresponses.Account, limit uint) ([]responses.Repository, responses.APIResponse) {
	namespaceAccountID := ""
	if acctPrefix != "" {
		acc, err := user.EnziSession.GetAccount(acctPrefix)
		if err != nil {
			apiErrs, _ := err.(*enzierrors.APIErrors)
			if dtrutil.CheckContainsEnziError(apiErrs, "NO_SUCH_ACCOUNT") {
				if failIfAccountNotExists {
					return nil, responses.APIError(errors.NoSuchAccountError(acctPrefix))
				}
			} else {
				return nil, responses.APIError(errors.InternalError(ctx, err))
			}
		} else {
			namespaceAccountID = acc.ID
		}
	}

	namespaceAccountIDToAccount := map[string]*enziresponses.Account{}
	namespaceAccountIDs := make([]string, len(accountResults))
	for i, acc := range accountResults {
		namespaceAccountIDToAccount[acc.ID] = &accountResults[i]
		namespaceAccountIDs[i] = acc.ID
	}
	namespaceMatchRepos, _, err := a.repoMgr.ListRepositoriesInNamespaces(namespaceAccountIDs, "", limit)
	if err != nil {
		return nil, responses.APIError(errors.InternalError(ctx, err))
	}

	nameMatchRepos, err := a.repoMgr.AutocompleteAllRepositories(repoPrefix, limit, namespaceAccountID)
	if err != nil {
		return nil, responses.APIError(errors.InternalError(ctx, err))
	}

	groupedRepos := groupReposByNamespace(concatRepoLists(namespaceMatchRepos, nameMatchRepos))
	visibleRepoIDSet := map[string]struct{}{}
	for namespaceAccountID, repos := range groupedRepos {
		if _, ok := namespaceAccountIDToAccount[namespaceAccountID]; !ok {
			ns, err := user.EnziSession.GetAccount("id:" + namespaceAccountID)
			if err != nil {
				apiErrs, _ := err.(*enzierrors.APIErrors)
				if dtrutil.CheckContainsEnziError(apiErrs, "NO_SUCH_ACCOUNT") {
					delete(groupedRepos, namespaceAccountID)
					// TODO remove account's repos from DTR???
					continue
				} else {
					return nil, responses.APIError(errors.InternalError(ctx, err))
				}
			}
			namespaceAccountIDToAccount[namespaceAccountID] = ns
		}
		ns := namespaceAccountIDToAccount[namespaceAccountID]

		accessLevel, err := a.authorizer.NamespaceAccess(user, ns)
		if err != nil {
			return nil, responses.APIError(errors.InternalError(ctx, err))
		}

		filteredRepos, err := a.authorizer.FilterVisibleReposInNamespace(user, ns, repos, accessLevel == schema.AccessLevelAdmin)
		if err != nil {
			return nil, responses.APIError(errors.InternalError(ctx, err))
		}

		for _, repo := range filteredRepos {
			visibleRepoIDSet[repo.ID] = struct{}{}
		}
	}

	var filteredNamespaceMatchRepos []*schema.Repository
	for _, repo := range namespaceMatchRepos {
		if _, ok := visibleRepoIDSet[repo.ID]; ok {
			filteredNamespaceMatchRepos = append(filteredNamespaceMatchRepos, repo)
			delete(visibleRepoIDSet, repo.ID) // for uniqueness in the result list
		}
	}
	var filteredNameMatchRepos []*schema.Repository
	for _, repo := range nameMatchRepos {
		if _, ok := visibleRepoIDSet[repo.ID]; ok {
			filteredNameMatchRepos = append(filteredNameMatchRepos, repo)
		}
	}

	numRepos := uint(len(visibleRepoIDSet))
	namespaceMatchReposLimit := limit / 4
	if numRepos > limit {
		numRepos = limit
		if uint(len(filteredNamespaceMatchRepos)) > namespaceMatchReposLimit &&
			uint(len(filteredNameMatchRepos))+namespaceMatchReposLimit >= limit {
			filteredNamespaceMatchRepos = filteredNamespaceMatchRepos[:namespaceMatchReposLimit]
		}
		filteredNameMatchRepos = filteredNameMatchRepos[:numRepos-uint(len(filteredNamespaceMatchRepos))]
	}

	reposResponse := make([]responses.Repository, 0, numRepos)
	for _, repo := range concatRepoLists(filteredNamespaceMatchRepos, filteredNameMatchRepos) {
		reposResponse = append(reposResponse, responses.MakeRepository(namespaceAccountIDToAccount[repo.NamespaceAccountID].Name, namespaceAccountIDToAccount[repo.NamespaceAccountID].IsOrg, repo, false))
	}
	return reposResponse, nil
}

func groupReposByNamespace(repos []*schema.Repository) map[string][]*schema.Repository {
	groups := make(map[string][]*schema.Repository)
	for _, repo := range repos {
		groups[repo.NamespaceAccountID] = append(groups[repo.NamespaceAccountID], repo)
	}
	return groups
}

func concatRepoLists(l1, l2 []*schema.Repository) []*schema.Repository {
	unionRepos := make([]*schema.Repository, 0, len(l1)+len(l2))
	unionRepos = append(unionRepos, l1...)
	unionRepos = append(unionRepos, l2...)
	return unionRepos
}

func parseSearchOptions(r *restful.Request) forms.SearchOptions {
	var opts forms.SearchOptions
	optionHandlers := map[string]func(string){
		"query":               func(s string) { opts.Query = s },
		"includeRepositories": func(b string) { opts.IncludeRepositories, _ = strconv.ParseBool(b) },
		"includeAccounts":     func(b string) { opts.IncludeAccounts, _ = strconv.ParseBool(b) },
		"namespace":           func(s string) { opts.Namespace = s },
		"limit":               func(i string) { limit, _ := strconv.ParseUint(i, 10, 64); opts.Limit = uint(limit) },
	}
	for option, handler := range optionHandlers {
		if v := r.QueryParameter(option); v != "" {
			handler(v)
		}
	}
	return opts
}
