'use strict';

var deployWizardTemplate = require('./deploy.wizard.html');
var deployYmlTemplate = require('./deploy.yml.html');
var deployStackTemplate = require('./deploy.stack.html');
require('./deploy.css');

getRoutes.$inject = ['$stateProvider', '$urlRouterProvider'];
function getRoutes($stateProvider, $urlRouterProvider) {
  $stateProvider
  .state('dashboard.deploywizard', {
    url: '^/deploy/wizard',
    template: deployWizardTemplate,
    controller: 'DeployWizardController',
    controllerAs: 'vm',
    authenticate: true,
    resolve: {
      permissions: ['AccountsService', '$state', 'AuthService', function (AccountsService, $state, AuthService) {
        if(AuthService.getCurrentUser().admin) {
          return AccountsService.getPermissions().then(null, function(errorData) {
            $state.go('error');
          });
        } else {
          return AccountsService.getPermissionsForAccount(AuthService.getCurrentUser().username).then(null, function(errorData) {
            $state.go('error');
          });
        }
      }],
      networks: ['NetworksService', '$state', function(NetworksService, $state) {
        return NetworksService.list().then(null, function(error) {
          // FIXME: This is a temporary check where we will be treating the 403 status code as 'empty',
          // since it's trickier on the backend to return an empty list of networks when the user doesn't
          // have the necessary permissions
          if(error.status !== 403) {
            $state.go('error');
          }
          return [];
        });
      }],
      volumes: ['VolumesService', '$state', function(VolumesService, $state) {
        return VolumesService.list().then(null, function(error) {
          $state.go('error');
        });
      }]
    },
    ncyBreadcrumb: {
      parent: 'dashboard.resources',
      label: 'Deploy'
    }
  })
  .state('dashboard.deploystack', {
    url: '^/deploy/stack',
    template: deployStackTemplate,
    authenticate: true,
    controller: 'DeployStackController',
    controllerAs: 'vm',
    ncyBreadcrumb: {
      parent: 'dashboard.resources',
      label: 'Deploy'
    }
  });
}

module.exports = getRoutes;
