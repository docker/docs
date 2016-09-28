package filters

import (
	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/api/server/internal/helpers"
	"github.com/docker/orca/enzi/authn"
	"github.com/docker/orca/enzi/authz"
	"github.com/docker/orca/enzi/config"
	"github.com/docker/orca/enzi/schema"
	"github.com/emicklei/go-restful"
)

// RequestData holds context about an API request and is able to gather and
// validate data required to handle the request. Any fields that are not
// needed are left empty or nil.
type RequestData struct {
	schemaMgr     schema.Manager
	authorizer    authz.Authorizer
	ctx           context.Context
	r             *restful.Request
	clientAccount *authn.Account // The authenticated client. May be anonymous.

	AuthConfig *config.Auth

	Acct    *schema.Account // An account in the context of the request.
	User    *schema.Account // The user in the context of the request.
	Org     *schema.Account // The organization in the context of the request.
	Team    *schema.Team    // The team in the context of the request.
	Service *schema.Service // The service in the context of the request.

	ClientAccountAccess *authz.MembershipAccess
	ClientOrgAccess     *authz.MembershipAccess
	ClientTeamAccess    *authz.MembershipAccess

	filters     []func()
	errResponse responses.APIResponse
}

// NewRequestData creates a new RequestData object which can gather and
// validate data for the given request and client account using the given
// schema manager and authorizer.
func NewRequestData(baseContext context.Context, schemaMgr schema.Manager, authorizer authz.Authorizer, r *restful.Request, clientAccount *authn.Account) *RequestData {
	return &RequestData{
		ctx:           baseContext,
		schemaMgr:     schemaMgr,
		authorizer:    authorizer,
		r:             r,
		clientAccount: clientAccount,
	}
}

// AddFilters appends the given functions to the filters list on this
// RequestData object. Each filter will be called in order to gather and
// validate data pertaining to the request. Usually, these filter functions are
// simply methods to call on this RequestData object.
func (rd *RequestData) AddFilters(filters ...func()) *RequestData {
	rd.filters = append(rd.filters, filters...)
	return rd
}

// EvaluateFilters evaluates all filters in this RequestData's filter list,
// stopping if any error is encountered. If any error is encountered, a non-nil
// APIResponse representing the error is returned; returns nil otherwise. Once
// called, the internal filters list is cleared. You must add more filters in
// order to perform further evaluation.
func (rd *RequestData) EvaluateFilters() (errResponse responses.APIResponse) {
	defer func() {
		// Clear all filters.
		rd.filters = nil
	}()

	for _, filter := range rd.filters {
		if filter(); rd.errResponse != nil {
			return rd.errResponse
		}
	}

	return nil
}

// GetAuthConfig gets the current auth configuration from the authorizer.
func (rd *RequestData) GetAuthConfig() {
	rd.AuthConfig, rd.errResponse = helpers.AuthConfig(rd.ctx, rd.authorizer)
}

// MakeFilterGetAccount returns a filter function which gets the account in
// question or sets a NotFound error response.
func (rd *RequestData) MakeFilterGetAccount(nameOrID string) func() {
	return func() {
		rd.Acct, rd.errResponse = helpers.AccountOrNotFound(rd.ctx, rd.schemaMgr, nameOrID)
	}
}

// MakeFilterGetUser returns a filter function which gets the user in question
// or sets a NotFound error response.
func (rd *RequestData) MakeFilterGetUser(nameOrID string) func() {
	return func() {
		rd.User, rd.errResponse = helpers.UserOrNotFound(rd.ctx, rd.schemaMgr, nameOrID)
	}
}

// MakeFilterGetOrganization returns a filter function which gets the
// organization in question or sets a NotFound error response.
func (rd *RequestData) MakeFilterGetOrganization(nameOrID string) func() {
	return func() {
		rd.Org, rd.errResponse = helpers.OrgOrNotFound(rd.ctx, rd.schemaMgr, nameOrID)
	}
}

// MakeFilterGetService returns a filter function which gets the service in
// question owned by the current account or sets a NotFound error.
func (rd *RequestData) MakeFilterGetService(nameOrID string) func() {
	return func() {
		rd.Service, rd.errResponse = helpers.ServiceOrNotFound(rd.ctx, rd.schemaMgr, rd.Acct.ID, nameOrID)
	}
}

// MakeFilterGetTeam returns a filter function which gets the team in question
// in the current organization or sets a NotFound error.
func (rd *RequestData) MakeFilterGetTeam(nameOrID string) func() {
	return func() {
		rd.Team, rd.errResponse = helpers.TeamOrNotFound(rd.ctx, rd.schemaMgr, rd.Org.ID, nameOrID)
	}
}

// GetAccountAccess gets the access that the client user has in the current
// account. If the account is an org, it will be retrieved as OrgAccess would.
// if the account is a user then account access will only be set if the client
// user *is* that account or if the client user is an admin.
func (rd *RequestData) GetAccountAccess() {
	rd.ClientAccountAccess, rd.errResponse = helpers.AccountAccess(rd.ctx, rd.authorizer, rd.Acct, rd.clientAccount)
}

// RequireAccountAdmin ensures that the client user's account access is not
// nil and that the IsAdmin field is true.
func (rd *RequestData) RequireAccountAdmin() {
	if rd.ClientAccountAccess == nil || !rd.ClientAccountAccess.IsAdmin {
		rd.errResponse = responses.APIError(errors.NotAuthorized(authz.ErrAccountAdminAccessRequired.Error()))
	}
}

// GetOrgAccess gets the membership access that the client user has in the
// current organization.
func (rd *RequestData) GetOrgAccess() {
	rd.ClientOrgAccess, rd.errResponse = helpers.OrgAccess(rd.ctx, rd.authorizer, rd.Org.ID, rd.clientAccount)
}

// RequireOrgMember ensures that the client user's orgAccess is not nil.
func (rd *RequestData) RequireOrgMember() {
	if rd.ClientOrgAccess == nil {
		rd.errResponse = responses.APIError(errors.NotAuthorized(authz.ErrOrgMemberAccessRequired.Error()))
	}
}

// RequireOrgAdmin ensures that the client user's orgAccess is not nil and that
// the IsAdmin field is true.
func (rd *RequestData) RequireOrgAdmin() {
	if rd.ClientOrgAccess == nil || !rd.ClientOrgAccess.IsAdmin {
		rd.errResponse = responses.APIError(errors.NotAuthorized(authz.ErrOrgAdminAccessRequired.Error()))
	}
}

// GetTeamAccess gets the membership access that the client user has in the
// current team.
func (rd *RequestData) GetTeamAccess() {
	rd.ClientTeamAccess, rd.errResponse = helpers.TeamAccess(rd.ctx, rd.authorizer, *rd.ClientOrgAccess, rd.Team.ID, rd.clientAccount)
}

// RequireTeamMember ensures that the client user's teamAccess is not nil.
func (rd *RequestData) RequireTeamMember() {
	if rd.ClientTeamAccess == nil {
		rd.errResponse = responses.APIError(errors.NotAuthorized(authz.ErrTeamMemberAccessRequired.Error()))
	}
}

// RequireTeamAdmin ensures that the client user's teamAccess is not nil and
// that the IsAdmin field is true.
func (rd *RequestData) RequireTeamAdmin() {
	if rd.ClientTeamAccess == nil || !rd.ClientTeamAccess.IsAdmin {
		rd.errResponse = responses.APIError(errors.NotAuthorized(authz.ErrTeamAdminAccessRequired.Error()))
	}
}
