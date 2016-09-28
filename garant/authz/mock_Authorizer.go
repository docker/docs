package authz

import "net/http"

import "github.com/stretchr/testify/mock"

import "github.com/docker/distribution/context"

import "github.com/docker/dhe-deploy/garant/authn"

import "github.com/docker/dhe-deploy/manager/schema"

import enziresponses "github.com/docker/orca/enzi/api/responses"

import "github.com/docker/garant/auth"

type MockAuthorizer struct {
	mock.Mock
}

// Authenticate provides a mock function with given fields: ctx, r
func (_m *MockAuthorizer) Authenticate(ctx context.Context, r *http.Request) (auth.Account, error) {
	ret := _m.Called(ctx, r)

	var r0 auth.Account
	if rf, ok := ret.Get(0).(func(context.Context, *http.Request) auth.Account); ok {
		r0 = rf(ctx, r)
	} else {
		r0 = ret.Get(0).(auth.Account)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *http.Request) error); ok {
		r1 = rf(ctx, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Authorize provides a mock function with given fields: ctx, acct, service, requestedAccess
func (_m *MockAuthorizer) Authorize(ctx context.Context, acct auth.Account, service string, requestedAccess ...auth.Access) ([]auth.Access, error) {
	ret := _m.Called(ctx, acct, service, requestedAccess)

	var r0 []auth.Access
	if rf, ok := ret.Get(0).(func(context.Context, auth.Account, string, ...auth.Access) []auth.Access); ok {
		r0 = rf(ctx, acct, service, requestedAccess...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]auth.Access)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, auth.Account, string, ...auth.Access) error); ok {
		r1 = rf(ctx, acct, service, requestedAccess...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthenticateWithToken provides a mock function with given fields: ctx, token
func (_m *MockAuthorizer) AuthenticateWithToken(ctx context.Context, token string) (auth.Account, error) {
	ret := _m.Called(ctx, token)

	var r0 auth.Account
	if rf, ok := ret.Get(0).(func(context.Context, string) auth.Account); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Get(0).(auth.Account)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthenticateWithPassword provides a mock function with given fields: ctx, username, password
func (_m *MockAuthorizer) AuthenticateWithPassword(ctx context.Context, username string, password string) (auth.Account, error) {
	ret := _m.Called(ctx, username, password)

	var r0 auth.Account
	if rf, ok := ret.Get(0).(func(context.Context, string, string) auth.Account); ok {
		r0 = rf(ctx, username, password)
	} else {
		r0 = ret.Get(0).(auth.Account)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetToken provides a mock function with given fields: ctx, acct, clientID
func (_m *MockAuthorizer) GetToken(ctx context.Context, acct auth.Account, clientID string) (auth.RefreshToken, error) {
	ret := _m.Called(ctx, acct, clientID)

	var r0 auth.RefreshToken
	if rf, ok := ret.Get(0).(func(context.Context, auth.Account, string) auth.RefreshToken); ok {
		r0 = rf(ctx, acct, clientID)
	} else {
		r0 = ret.Get(0).(auth.RefreshToken)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, auth.Account, string) error); ok {
		r1 = rf(ctx, acct, clientID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AccountTokens provides a mock function with given fields: ctx, acct
func (_m *MockAuthorizer) AccountTokens(ctx context.Context, acct auth.Account) ([]auth.RefreshToken, error) {
	ret := _m.Called(ctx, acct)

	var r0 []auth.RefreshToken
	if rf, ok := ret.Get(0).(func(context.Context, auth.Account) []auth.RefreshToken); ok {
		r0 = rf(ctx, acct)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]auth.RefreshToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, auth.Account) error); ok {
		r1 = rf(ctx, acct)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RevokeToken provides a mock function with given fields: ctx, token
func (_m *MockAuthorizer) RevokeToken(ctx context.Context, token string) error {
	ret := _m.Called(ctx, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthenticateRequestUser provides a mock function with given fields: ctx, r
func (_m *MockAuthorizer) AuthenticateRequestUser(ctx context.Context, r *http.Request) (*authn.User, error) {
	ret := _m.Called(ctx, r)

	var r0 *authn.User
	if rf, ok := ret.Get(0).(func(context.Context, *http.Request) *authn.User); ok {
		r0 = rf(ctx, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*authn.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *http.Request) error); ok {
		r1 = rf(ctx, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckAdminOrOrgOwner provides a mock function with given fields: user, orgID
func (_m *MockAuthorizer) CheckAdminOrOrgOwner(user *authn.User, orgID string) error {
	ret := _m.Called(user, orgID)

	var r0 error
	if rf, ok := ret.Get(0).(func(*authn.User, string) error); ok {
		r0 = rf(user, orgID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CheckAdminOrOrgMember provides a mock function with given fields: user, orgID
func (_m *MockAuthorizer) CheckAdminOrOrgMember(user *authn.User, orgID string) error {
	ret := _m.Called(user, orgID)

	var r0 error
	if rf, ok := ret.Get(0).(func(*authn.User, string) error); ok {
		r0 = rf(user, orgID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NamespaceAccess provides a mock function with given fields: user, ns
func (_m *MockAuthorizer) NamespaceAccess(user *authn.User, ns *enziresponses.Account) (string, error) {
	ret := _m.Called(user, ns)

	var r0 string
	if rf, ok := ret.Get(0).(func(*authn.User, *enziresponses.Account) string); ok {
		r0 = rf(user, ns)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*authn.User, *enziresponses.Account) error); ok {
		r1 = rf(user, ns)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RepositoryAccess provides a mock function with given fields: user, repo, ns
func (_m *MockAuthorizer) RepositoryAccess(user *authn.User, repo *schema.Repository, ns *enziresponses.Account) (string, error) {
	ret := _m.Called(user, repo, ns)

	var r0 string
	if rf, ok := ret.Get(0).(func(*authn.User, *schema.Repository, *enziresponses.Account) string); ok {
		r0 = rf(user, repo, ns)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*authn.User, *schema.Repository, *enziresponses.Account) error); ok {
		r1 = rf(user, repo, ns)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AllVisibleRepositoriesInNamespace provides a mock function with given fields: user, ns
func (_m *MockAuthorizer) AllVisibleRepositoriesInNamespace(user *authn.User, ns *enziresponses.Account) ([]*schema.Repository, error) {
	ret := _m.Called(user, ns)

	var r0 []*schema.Repository
	if rf, ok := ret.Get(0).(func(*authn.User, *enziresponses.Account) []*schema.Repository); ok {
		r0 = rf(user, ns)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*schema.Repository)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*authn.User, *enziresponses.Account) error); ok {
		r1 = rf(user, ns)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VisibleRepositoriesInNamespace provides a mock function with given fields: user, ns, startPK, limit
func (_m *MockAuthorizer) VisibleRepositoriesInNamespace(user *authn.User, ns *enziresponses.Account, startPK string, limit uint) ([]*schema.Repository, string, error) {
	ret := _m.Called(user, ns, startPK, limit)

	var r0 []*schema.Repository
	if rf, ok := ret.Get(0).(func(*authn.User, *enziresponses.Account, string, uint) []*schema.Repository); ok {
		r0 = rf(user, ns, startPK, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*schema.Repository)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(*authn.User, *enziresponses.Account, string, uint) string); ok {
		r1 = rf(user, ns, startPK, limit)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*authn.User, *enziresponses.Account, string, uint) error); ok {
		r2 = rf(user, ns, startPK, limit)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// FilterVisibleReposInNamespace provides a mock function with given fields: user, namespace, repos, isNamespaceAdmin
func (_m *MockAuthorizer) FilterVisibleReposInNamespace(user *authn.User, namespace *enziresponses.Account, repos []*schema.Repository, isNamespaceAdmin bool) ([]*schema.Repository, error) {
	ret := _m.Called(user, namespace, repos, isNamespaceAdmin)

	var r0 []*schema.Repository
	if rf, ok := ret.Get(0).(func(*authn.User, *enziresponses.Account, []*schema.Repository, bool) []*schema.Repository); ok {
		r0 = rf(user, namespace, repos, isNamespaceAdmin)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*schema.Repository)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*authn.User, *enziresponses.Account, []*schema.Repository, bool) error); ok {
		r1 = rf(user, namespace, repos, isNamespaceAdmin)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SharedRepositoriesForUser provides a mock function with given fields: user, startPK, limit
func (_m *MockAuthorizer) SharedRepositoriesForUser(user *authn.User, startPK string, limit uint) ([]*schema.Repository, string, error) {
	ret := _m.Called(user, startPK, limit)

	var r0 []*schema.Repository
	if rf, ok := ret.Get(0).(func(*authn.User, string, uint) []*schema.Repository); ok {
		r0 = rf(user, startPK, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*schema.Repository)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(*authn.User, string, uint) string); ok {
		r1 = rf(user, startPK, limit)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*authn.User, string, uint) error); ok {
		r2 = rf(user, startPK, limit)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
