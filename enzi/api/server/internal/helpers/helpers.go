package helpers

import (
	"math"
	"strconv"
	"strings"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/authn"
	ldapconfig "github.com/docker/orca/enzi/authn/ldap/config"
	"github.com/docker/orca/enzi/authz"
	"github.com/docker/orca/enzi/config"
	"github.com/docker/orca/enzi/schema"
	"github.com/emicklei/go-restful"
)

// PageParams gets the page start offset and page results limit from the given
// request using query parameters with the given names.
func PageParams(r *restful.Request, startParamName, limitParamName string) (start string, limit uint) {
	start = r.QueryParameter(startParamName)
	limitStr := r.QueryParameter(limitParamName)

	parsedLimit, _ := strconv.ParseUint(limitStr, 10, 32)
	if parsedLimit == 0 {
		parsedLimit = api.DefaultPerPageLimit
	}

	limit = uint(parsedLimit)
	if limit > api.MaxPerPageLimit {
		limit = api.MaxPerPageLimit
	}

	return start, limit
}

// ParseOffsetAsUint is used for paging when the offset should be a numbered
// offset rather than a string identifie
func ParseOffsetAsUint(offsetStr string) uint {
	offset, _ := strconv.ParseUint(offsetStr, 10, 32)
	if offset > math.MaxUint32 {
		offset = math.MaxUint32
	}

	return uint(offset)
}

// ParseBoolQueryParam attempts to parse a query parameter with the given name
// into boolean. If there is an error or the query parameter is not given, the
// default value is used.
func ParseBoolQueryParam(r *restful.Request, queryParamName string, defaultVal bool) bool {
	val, err := strconv.ParseBool(r.QueryParameter(queryParamName))
	if err != nil {
		return defaultVal
	}

	return val
}

// AuthConfig retrieves the current auth configuration. Returns a non-nil
// responses.APIResponse describing any error.
func AuthConfig(ctx context.Context, authorizer authz.Authorizer) (authConfig *config.Auth, errResponse responses.APIResponse) {
	authConfig, err := authorizer.AuthConfig()
	if err != nil {
		return nil, responses.APIError(errors.Internal(ctx, err))
	}

	return authConfig, nil
}

// LDAPSettings retrieves the current LDAP configuration. Returns a non-nil
// responses.APIResponse describing any error.
func LDAPSettings(ctx context.Context, authorizer authz.Authorizer) (ldapSettings *ldapconfig.Settings, errResponse responses.APIResponse) {
	ldapSettings, err := authorizer.LDAPSettings()
	if err != nil {
		return nil, responses.APIError(errors.Internal(ctx, err))
	}

	return ldapSettings, nil
}

// OpenIDConfig retrieves the current system OpenID configuration. Returns
// a non-nil responses.APIResponse describing any error.
func OpenIDConfig(ctx context.Context, mgr schema.Manager) (openIDConfig *config.OpenID, errResponse responses.APIResponse) {
	openIDConfig, err := config.GetOpenIDConfig(mgr)
	if err != nil {
		return nil, responses.APIError(errors.Internal(ctx, err))
	}

	return openIDConfig, nil
}

const idSpecifierPrefix = "id:"

// AccountByNameOrID retrieves the account with the given name or ID using the
// given schema manager. Handles parsing the input as an account name or ID.
func AccountByNameOrID(mgr schema.Manager, nameOrID string) (user *schema.Account, err error) {
	if strings.HasPrefix(nameOrID, idSpecifierPrefix) {
		accountID := strings.TrimPrefix(nameOrID, idSpecifierPrefix)
		return mgr.GetAccountByID(accountID)
	}

	return mgr.GetAccountByName(forms.NormalizeAccountName(nameOrID))
}

// getAccountHelper generalizes the proces of getting an account (or more
// specifically a user or organization) by using the given ifName or ifID
// functions when the given nameOrID is either an account name or account ID,
// respectively. If there is no such account or some other error occurs then a
// non-nil responses.APIResponse is returned indicating the error with
// appropriate status code. Returns a non-nil account on success.
func getAccountHelper(ctx context.Context, nameOrID string, ifName, ifID func(string) (*schema.Account, error)) (acct *schema.Account, errResponse responses.APIResponse) {
	var err error

	if strings.HasPrefix(nameOrID, idSpecifierPrefix) {
		acctID := strings.TrimPrefix(nameOrID, idSpecifierPrefix)
		acct, err = ifID(acctID)
	} else {
		acct, err = ifName(forms.NormalizeAccountName(nameOrID))
	}

	switch err {
	case nil:
		return acct, nil
	case schema.ErrNoSuchAccount:
		return nil, responses.APIError(errors.NoSuchAccount(nameOrID))
	default:
		return nil, responses.APIError(errors.Internal(ctx, err))
	}
}

// GetAccountByNameOrID determines if the given nameOrID is in fact a name or
// ID and retrieves the corresponding account from the backend.
func GetAccountByNameOrID(mgr schema.Manager, nameOrID string) (*schema.Account, error) {
	if strings.HasPrefix(nameOrID, idSpecifierPrefix) {
		acctID := strings.TrimPrefix(nameOrID, idSpecifierPrefix)
		return mgr.GetAccountByID(acctID)
	}

	return mgr.GetAccountByName(forms.NormalizeAccountName(nameOrID))
}

// AccountOrNotFound attempts to retrieve the account with the given name or
// ID. If there is no such account or some other error occurs then a non-nil
// responses.APIResponse is returned indicating the error with appropriate
// status code. Returns a non-nil account on success.
func AccountOrNotFound(ctx context.Context, mgr schema.Manager, nameOrID string) (user *schema.Account, errResponse responses.APIResponse) {
	return getAccountHelper(ctx, nameOrID, mgr.GetAccountByName, mgr.GetAccountByID)
}

// UserOrNotFound attempts to retrieve the user account with the given name
// or ID. If there is no such account or some other error occurs then a non-nil
// responses.APIResponse is returned indicating the error with appropriate
// status code. Returns a non-nil account on success.
func UserOrNotFound(ctx context.Context, mgr schema.Manager, nameOrID string) (user *schema.Account, errResponse responses.APIResponse) {
	return getAccountHelper(ctx, nameOrID, mgr.GetUserByName, mgr.GetUserByID)
}

// OrgOrNotFound attempts to retrieve the organization account with the given
// name or ID. If there is no such account or some other error occurs then a
// non-nil responses.APIResponse is returned indicating the error with
// appropriate status code. Returns a non-nil account on success.
func OrgOrNotFound(ctx context.Context, mgr schema.Manager, nameOrID string) (org *schema.Account, errResponse responses.APIResponse) {
	return getAccountHelper(ctx, nameOrID, mgr.GetOrgByName, mgr.GetOrgByID)
}

// AccountAccess retrieves the level of access that the given client account
// has for the given account. May return (nil, nil) indicating that the client
// account is not related to the account.
func AccountAccess(ctx context.Context, authorizer authz.Authorizer, acct *schema.Account, clientAccount *authn.Account) (acctAccess *authz.MembershipAccess, errResponse responses.APIResponse) {
	acctAccess, err := authorizer.AccountAccess(acct, clientAccount)
	if err != nil {
		return nil, responses.APIError(errors.Internal(ctx, err))
	}

	return acctAccess, nil
}

// OrgAccess retrieves the level of access that the given client Account has in
// the given organization. May return (nil, nil) indicating that the client
// user is not a member of the organization.
func OrgAccess(ctx context.Context, authorizer authz.Authorizer, orgID string, clientAccount *authn.Account) (orgAccess *authz.MembershipAccess, errResponse responses.APIResponse) {
	orgAccess, err := authorizer.OrgMembershipAccess(orgID, clientAccount)
	if err != nil {
		return nil, responses.APIError(errors.Internal(ctx, err))
	}

	return orgAccess, nil
}

// ServiceByNameOrID retrieves the service with the given name or ID using the
// given schema manager. Handles parsing the input as a service name or ID.
func ServiceByNameOrID(mgr schema.Manager, ownerID, nameOrID string) (service *schema.Service, err error) {
	if strings.HasPrefix(nameOrID, idSpecifierPrefix) {
		serviceID := strings.TrimPrefix(nameOrID, idSpecifierPrefix)
		service, err = mgr.GetServiceByID(serviceID)
		if err == nil && service.OwnerID != ownerID {
			err = schema.ErrNoSuchService
		}
	} else {
		nameOrID = forms.NormalizeFullName(nameOrID)
		service, err = mgr.GetServiceByName(ownerID, nameOrID)
	}

	return service, err
}

// ServiceOrNotFound attempts to retrieve the service with the given name or
// ID and has the given owner ID. If there is no such service or some other
// error occurs then a non-nil responses.APIResponse is returned indicating
// the error with appropriate status code. Returns a non-nil service on
// success.
func ServiceOrNotFound(ctx context.Context, mgr schema.Manager, ownerID, nameOrID string) (service *schema.Service, errResponse responses.APIResponse) {
	service, err := ServiceByNameOrID(mgr, ownerID, nameOrID)

	switch err {
	case nil:
		return service, nil
	case schema.ErrNoSuchService:
		return nil, responses.APIError(errors.NoSuchService(nameOrID))
	default:
		return nil, responses.APIError(errors.Internal(ctx, err))
	}

}

// TeamByNameOrID retrieves the team with the given name or ID using the given
// schema manager. Handles parsing the input as a team name or ID.
func TeamByNameOrID(mgr schema.Manager, orgID, nameOrID string) (team *schema.Team, err error) {
	if strings.HasPrefix(nameOrID, idSpecifierPrefix) {
		teamID := strings.TrimPrefix(nameOrID, idSpecifierPrefix)
		team, err = mgr.GetTeamByID(teamID)
		if err == nil && team.OrgID != orgID {
			err = schema.ErrNoSuchTeam
		}
	} else {
		nameOrID = forms.NormalizeFullName(nameOrID)
		team, err = mgr.GetTeamByName(orgID, nameOrID)
	}

	return team, err
}

// TeamOrNotFound attempts to retrieve the team with the given name or ID. If
// there is no such team or some other error occurs then a non-nil
// responses.APIResponse is returned indicating the error with appropriate
// status code. Returns a non-nil team on success.
func TeamOrNotFound(ctx context.Context, mgr schema.Manager, orgID, nameOrID string) (team *schema.Team, errResponse responses.APIResponse) {
	team, err := TeamByNameOrID(mgr, orgID, nameOrID)

	switch err {
	case nil:
		return team, nil
	case schema.ErrNoSuchTeam:
		return nil, responses.APIError(errors.NoSuchTeam(nameOrID))
	default:
		return nil, responses.APIError(errors.Internal(ctx, err))
	}

}

// TeamAccess retrieves the level of access that the given client account has
// in the given organization and team. May return (nil, nil, nil) indicating
// that the client account is not a member of the organization and team.
func TeamAccess(ctx context.Context, authorizer authz.Authorizer, orgAccess authz.MembershipAccess, teamID string, clientAccount *authn.Account) (teamAccess *authz.MembershipAccess, errResponse responses.APIResponse) {
	teamAccess, err := authorizer.TeamMembershipAccess(teamID, orgAccess, clientAccount)
	if err != nil {
		return nil, responses.APIError(errors.Internal(ctx, err))
	}

	return teamAccess, nil
}
