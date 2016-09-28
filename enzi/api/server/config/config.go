package config

import (
	"fmt"
	"net/http"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/api/server"
	"github.com/docker/orca/enzi/api/server/internal/helpers"
	"github.com/docker/orca/enzi/authn"
	"github.com/docker/orca/enzi/authn/ldap"
	ldapconfig "github.com/docker/orca/enzi/authn/ldap/config"
	"github.com/docker/orca/enzi/authz"
	"github.com/docker/orca/enzi/config"
	"github.com/docker/orca/enzi/passwords"
	"github.com/docker/orca/enzi/schema"
	"github.com/emicklei/go-restful"
)

// Service handles various API requests relating to auth config management.
type Service struct {
	server.Service

	schemaMgr  schema.Manager
	authorizer authz.Authorizer
}

// NewService returns a new Sessions Service.
func NewService(baseContext context.Context, schemaMgr schema.Manager, rootPath string) *Service {
	service := &Service{
		Service: server.Service{
			WebService:  new(restful.WebService),
			BaseContext: baseContext,
		},
		schemaMgr:  schemaMgr,
		authorizer: authz.NewAuthorizer(schemaMgr),
	}

	service.connectRoutes(rootPath)

	return service
}

// connectRoutes registers all API endpoints on this service with paths
// relative to the given rootPath.
func (s *Service) connectRoutes(rootPath string) {
	s.WebService.Path(rootPath).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("Config")

	routes := []server.Route{
		s.routeGetAuthConfig(),
		s.routeSetAuthConfig(),
		s.routeGetLDAPSettings(),
		s.routeSetLDAPSettings(),
		s.routeLDAPTryLogin(),
		s.routeGetOpenIDConfig(),
		s.routeSetOpenIDConfig(),
	}

	for _, route := range routes {
		route.Register(&s.Service)
	}
}

// routeGetAuthConfig returns a route describing the GetAuthConfig endpoint.
func (s *Service) routeGetAuthConfig() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/auth",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleGetAuthConfig),
		Doc:     "Retrieve current system auth configuration",
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, current auth config returned.",
				Sample:  responses.AuthConfig{},
			},
		},
	}
}

// handleGetAuthConfig handles a request for getting the current system auth
// configuration. The endpoint does not require admin access as it only returns
// non-sensitive data about which auth backend is currently used. This
// information is useful for determining what controls to show in a user
// interface.
func (s *Service) handleGetAuthConfig(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	authConfig, errResponse := helpers.AuthConfig(ctx, s.authorizer)
	if errResponse != nil {
		return errResponse
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeAuthConfig(authConfig))
}

// routeSetAuthConfig returns a route describing the SetAuthConfig endpoint.
func (s *Service) routeSetAuthConfig() server.Route {
	return server.Route{
		Method:     "PUT",
		Path:       "/auth",
		Handler:    server.WrapHandlerWithAdminAccount(s.authorizer, s.handleSetAuthConfig),
		Doc:        "Set system auth configuration",
		BodySample: forms.AuthConfig{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, current auth config set.",
				Sample:  responses.AuthConfig{},
			},
		},
	}
}

// handleSetAuthConfig handles a request for setting the system auth
// configuration.
func (s *Service) handleSetAuthConfig(ctx context.Context, adminAccount *authn.Account, r *restful.Request) responses.APIResponse {
	defer r.Request.Body.Close()

	// Decode and validate the form.
	form := new(forms.AuthConfig)
	if formErrors := form.ValidateJSON(r.Request.Body); len(formErrors) > 0 {
		return responses.APIError(formErrors...)
	}

	authConfig := &config.Auth{
		Backend: form.Backend,
	}

	// In general, we should validate the backend config before switching
	// to that backend.
	if authConfig.Backend == config.AuthBackendLDAP {
		if errResponse := s.ensureRecoveryAdminUserConfigured(ctx); errResponse != nil {
			return errResponse
		}
	}

	if err := config.SetAuthConfig(s.schemaMgr, authConfig); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeAuthConfig(authConfig))
}

func (s *Service) ensureRecoveryAdminUserConfigured(ctx context.Context) (errResponse responses.APIResponse) {
	// Get the LDAP settings to confirm that the recovery admin user
	// exists. Don't bother validating anything else. If it's misconfigured
	// then the recovery user should still be able to fix it.
	ldapSettings, errResponse := helpers.LDAPSettings(ctx, s.authorizer)
	if errResponse != nil {
		return errResponse
	}

	recoveryAdminUser, err := s.schemaMgr.GetUserByName(ldapSettings.RecoveryAdminUsername)
	if err != nil {
		if err == schema.ErrNoSuchAccount {
			// Don't allow the client to set the system to use the
			// LDAP auth backend if the recovery admin user does
			// not exist.
			return responses.APIError(errors.InvalidFormField("backend", "The recovery admin user must be created to use the LDAP auth backend"))
		}

		return responses.APIError(errors.Internal(ctx, err))
	}

	// Ensure they have a password hash.
	if recoveryAdminUser.PasswordHash == "" {
		return responses.APIError(errors.InvalidFormField("backend", "The recovery admin user must have a password configured to use the LDAP auth backend"))
	}

	return nil
}

// routeGetLDAPSettings returns a route describing the GetLDAPSettings endpoint.
func (s *Service) routeGetLDAPSettings() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/auth/ldap",
		Handler: server.WrapHandlerWithAdminAccount(s.authorizer, s.handleGetLDAPSettings),
		Doc:     "Retrieve current system LDAP configuration",
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, current LDAP config returned.",
				Sample:  responses.LDAPSettings{},
			},
		},
	}
}

// handleGetLDAPSettings handles a request for getting the current system LDAP
// configuration.
func (s *Service) handleGetLDAPSettings(ctx context.Context, adminAccount *authn.Account, r *restful.Request) responses.APIResponse {
	ldapSettings, errResponse := helpers.LDAPSettings(ctx, s.authorizer)
	if errResponse != nil {
		return errResponse
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeLdapSettings(ldapSettings))
}

// routeSetLDAPSettings returns a route describing the SetLDAPSettings endpoint.
func (s *Service) routeSetLDAPSettings() server.Route {
	return server.Route{
		Method:     "PUT",
		Path:       "/auth/ldap",
		Handler:    server.WrapHandlerWithAdminAccount(s.authorizer, s.handleSetLDAPSettings),
		Doc:        "Set system LDAP configuration",
		BodySample: forms.LDAPSettings{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, current LDAP config set.",
				Sample:  responses.LDAPSettings{},
			},
		},
	}
}

// handleSetLDAPSettings handles a request for setting the system LDAP
// configuration.
func (s *Service) handleSetLDAPSettings(ctx context.Context, adminAccount *authn.Account, r *restful.Request) responses.APIResponse {
	defer r.Request.Body.Close()

	// Decode and validate the form.
	form := new(forms.LDAPSettings)
	if formErrors := form.ValidateJSON(r.Request.Body); len(formErrors) > 0 {
		return responses.APIError(formErrors...)
	}

	// Perform further validaiton for the recover admin user settings.
	if errResponse := s.configureLDAPRecoveryAdminUser(ctx, form); errResponse != nil {
		return errResponse
	}

	ldapSettings := makeLDAPSettings(form)
	if err := ldapconfig.SetLDAPConfig(s.schemaMgr, ldapSettings); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeLdapSettings(ldapSettings))
}

func (s *Service) configureLDAPRecoveryAdminUser(ctx context.Context, form *forms.LDAPSettings) (errResponse responses.APIResponse) {
	var (
		passwordHash string
		err          error
	)

	// Hash the password if it was given.
	if form.RecoveryAdminPassword != nil {
		passwordHash, err = passwords.HashPassword(*form.RecoveryAdminPassword)
		if err != nil {
			return responses.APIError(errors.Internal(ctx, err))
		}
	}

	// Create the recovery admin user if it does not already exist. The
	// form's validation will have already normalized the account name.
	recoveryAdminUser, err := s.schemaMgr.GetUserByName(form.RecoveryAdminUsername)
	if err != nil {
		if err != schema.ErrNoSuchAccount {
			return responses.APIError(errors.Internal(ctx, err))
		}

		// The user does not yet exist. Ensure the password was given.
		if passwordHash == "" {
			return responses.APIError(errors.InvalidFormField("recoveryAdminPassword", "This field can not be left blank if the recovery admin user does not yet exist"))
		}

		recoveryAdminUser = &schema.Account{
			Name:    form.RecoveryAdminUsername,
			IsAdmin: true, IsActive: true,
			PasswordHash: passwordHash,
		}

		if err := s.schemaMgr.CreateAccount(recoveryAdminUser); err != nil {
			// Could be account already exists but it's an acceptable race condition for now.
			return responses.APIError(errors.Internal(ctx, err))
		}

		return nil
	}

	// The account exists - ensure it is a user.
	if recoveryAdminUser.IsOrg {
		return responses.APIError(errors.AccountExists())
	}

	// If the user doesn't already have a password hash, a password
	// must be provided.
	if recoveryAdminUser.PasswordHash == "" && passwordHash == "" {
		return responses.APIError(errors.InvalidFormField("recoveryAdminPassword", "This field can not be left blank if the recovery admin user does not yet have a local password hash"))
	}

	// If a new password was provided, we need to set it.
	if passwordHash != "" {
		updateFields := schema.AccountUpdateFields{
			PasswordHash: &passwordHash,
		}

		if err := s.schemaMgr.UpdateAccount(recoveryAdminUser.ID, updateFields); err != nil {
			return responses.APIError(errors.Internal(ctx, err))
		}
	}

	return nil
}

// makeLDAPSettings converts the given form (which should be pre-validated) to
// a backend ldap settings value.
func makeLDAPSettings(form *forms.LDAPSettings) *ldapconfig.Settings {
	userSearchConfigs := make([]ldapconfig.UserSearchOpts, len(form.UserSearchConfigs))
	for i, userSearchConfig := range form.UserSearchConfigs {
		userSearchConfigs[i] = ldapconfig.UserSearchOpts{
			BaseDN:       userSearchConfig.BaseDN,
			ScopeSubtree: userSearchConfig.ScopeSubtree,
			UsernameAttr: userSearchConfig.UsernameAttr,
			FullNameAttr: userSearchConfig.FullNameAttr,
			Filter:       userSearchConfig.Filter,
		}
	}

	return &ldapconfig.Settings{
		RecoveryAdminUsername: form.RecoveryAdminUsername,
		ServerURL:             form.ServerURL,
		NoSimplePagination:    form.NoSimplePagination,
		StartTLS:              form.StartTLS,
		RootCerts:             form.RootCerts,
		TLSSkipVerify:         form.TLSSkipVerify,
		ReaderDN:              form.ReaderDN,
		ReaderPassword:        form.ReaderPassword,
		UserSearchConfigs:     userSearchConfigs,
		AdminSyncOpts: ldapconfig.MemberSyncOpts{
			EnableSync:         form.AdminSyncOpts.EnableSync,
			SelectGroupMembers: form.AdminSyncOpts.SelectGroupMembers,
			GroupDN:            form.AdminSyncOpts.GroupDN,
			GroupMemberAttr:    form.AdminSyncOpts.GroupMemberAttr,
			SearchBaseDN:       form.AdminSyncOpts.SearchBaseDN,
			SearchScopeSubtree: form.AdminSyncOpts.SearchScopeSubtree,
			SearchFilter:       form.AdminSyncOpts.SearchFilter,
		},
		SyncSchedule: form.SyncSchedule,
	}
}

// routeLDAPTryLogin returns a route describing the LDAPTryLogin endpoint.
func (s *Service) routeLDAPTryLogin() server.Route {
	return server.Route{
		Method:     "POST",
		Path:       "/auth/ldap/tryLogin",
		Handler:    server.WrapHandlerWithAdminAccount(s.authorizer, s.handleLDAPTryLogin),
		Doc:        "Attempt to login with the given LDAP settings",
		BodySample: forms.TryLdapLogin{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, authenticated user returned.",
				Sample:  responses.Account{},
			},
		},
	}
}

// handleLDAPTryLogin handles a request to attempt login using the given LDAP
// settings.
func (s *Service) handleLDAPTryLogin(ctx context.Context, adminAccount *authn.Account, r *restful.Request) responses.APIResponse {
	defer r.Request.Body.Close()

	// Decode and validate the form.
	form := new(forms.TryLdapLogin)
	if formErrors := form.ValidateJSON(r.Request.Body); len(formErrors) > 0 {
		return responses.APIError(formErrors...)
	}

	ldapSettings := makeLDAPSettings(&form.LDAPSettings)

	user, logBuf, err := ldap.CheckLoginWithSettings(form.Username, form.Password, ldapSettings)
	if err != nil {
		return responses.APIError(errors.InvalidFormField("authentication", fmt.Sprintf("%s - Logs: %s", err, logBuf.String())))
	}

	context.GetLogger(ctx).Info(logBuf.String())

	return responses.JSONResponse(http.StatusOK, responses.MakeAccount(&user.Account))
}

// routeGetOpenIDConfig returns a route describing the GetOpenIDConfig
// endpoint.
func (s *Service) routeGetOpenIDConfig() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/openid",
		Handler: server.WrapHandlerWithAdminAccount(s.authorizer, s.handleGetOpenIDConfig),
		Doc:     "Retrieve current system OpenID configuration",
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, current OpenID config returned.",
				Sample:  responses.OpenIDConfig{},
			},
		},
	}
}

// handleGetOpenIDConfig handles a request for getting the current system
// OpenID configuration. The endpoint requires admin access.
func (s *Service) handleGetOpenIDConfig(ctx context.Context, adminAccount *authn.Account, r *restful.Request) responses.APIResponse {
	openIDConfig, errResponse := helpers.OpenIDConfig(ctx, s.schemaMgr)
	if errResponse != nil {
		return errResponse
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeOpenIDConfig(openIDConfig))
}

// routeSetOpenIDConfig returns a route describing the SetOpenIDConfig
// endpoint.
func (s *Service) routeSetOpenIDConfig() server.Route {
	return server.Route{
		Method:     "PUT",
		Path:       "/openid",
		Handler:    server.WrapHandlerWithAdminAccount(s.authorizer, s.handleSetOpenIDConfig),
		Doc:        "Set system OpenID configuration",
		BodySample: forms.OpenIDConfig{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, current OpenID config set.",
				Sample:  responses.OpenIDConfig{},
			},
		},
	}
}

// handleSetOpenIDConfig handles a request for setting the system OpenID
// configuration.
func (s *Service) handleSetOpenIDConfig(ctx context.Context, adminAccount *authn.Account, r *restful.Request) responses.APIResponse {
	defer r.Request.Body.Close()

	// Decode and validate the form.
	form := new(forms.OpenIDConfig)
	if formErrors := form.ValidateJSON(r.Request.Body); len(formErrors) > 0 {
		return responses.APIError(formErrors...)
	}

	openIDConfig := &config.OpenID{
		IssuerIdentifier: form.IssuerIdentifier,
	}

	if err := config.SetOpenIDConfig(s.schemaMgr, openIDConfig); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeOpenIDConfig(openIDConfig))
}
