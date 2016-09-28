package accounts

import (
	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/server"
	"github.com/docker/orca/enzi/authz"
	"github.com/docker/orca/enzi/schema"
	"github.com/emicklei/go-restful"
)

// Service handles various API requests relating to account management,
// including managing Teams and membership.
type Service struct {
	server.Service

	schemaMgr  schema.Manager
	authorizer authz.Authorizer
}

// NewService returns a new Accounts Service.
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
		Doc("Accounts")

	routes := []server.Route{
		s.routeCreateAccounts(),
		s.routeListAccounts(),
		s.routeGetAccount(),
		s.routeUpdateAccount(),
		s.routeChangePassword(),
		s.routeDeleteAccount(),
		s.routeListUserOrganizations(),
		s.routeListOrganizationMembers(),
		s.routeListOrganizationPublicMembers(),
		s.routeGetOrganizationAdminSyncConfig(),
		s.routeSetOrganizationAdminSyncConfig(),
		s.routeAddOrganizationMember(),
		s.routeUpdateOrganizationMembership(),
		s.routeGetOrganizationMembership(),
		s.routeDeleteOrganizationMember(),
		s.routeCreateTeam(),
		s.routeListTeams(),
		s.routeGetTeam(),
		s.routeUpdateTeam(),
		s.routeDeleteTeam(),
		s.routeListOrganizationMemberTeams(),
		s.routeListTeamMembers(),
		s.routeListTeamPublicMembers(),
		s.routeGetTeamMemberSyncConfig(),
		s.routeSetTeamMemberSyncConfig(),
		s.routeAddTeamMember(),
		s.routeUpdateTeamMembership(),
		s.routeGetTeamMembership(),
		s.routeDeleteTeamMember(),
		s.routeCreateService(),
		s.routeListServices(),
		s.routeGetService(),
		s.routeUpdateService(),
		s.routeDeleteService(),
	}

	for _, route := range routes {
		route.Register(&s.Service)
	}
}
