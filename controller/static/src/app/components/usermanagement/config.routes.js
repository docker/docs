'use strict';

var baseTemplate = require('./accountsBase.html');
var usersTemplate = require('./users.html');
var teamTemplate = require('./team.html');
var teamPermissionsTemplate = require('./teamPermissions.html');
var teamSettingsTemplate = require('./teamSettings.html');

require('./accountsBase.css');

getRoutes.$inject = ['$stateProvider', '$urlRouterProvider'];
function getRoutes($stateProvider, $urlRouterProvider) {
  $stateProvider
    .state('dashboard.accounts', {
      redirectTo: 'dashboard.accounts.users',
      template: baseTemplate,
      controller: 'AccountsBaseController as vm',
      authenticate: true,
      resolve: {
        teams: ['AccountsService', '$state', function (AccountsService, $state) {
          return AccountsService.listTeams().then(null, function(errorData) {
            $state.go(errorData);
          });
        }]
      },
      ncyBreadcrumb: {
        label: 'User Management'
      }
    })
    .state('dashboard.accounts.team', {
      url: '^/teams/{teamId}',
      template: teamTemplate,
      controller: 'TeamController as vm',
      authenticate: true,
      resolve: {
        team: ['AccountsService', '$state', '$stateParams', function (AccountsService, $state, $stateParams) {
          return AccountsService.getTeam($stateParams.teamId).then(null, function(errorData) {
            $state.go('error');
          });
        }],
        teamMembers: ['AccountsService', '$state', '$stateParams', function (AccountsService, $state, $stateParams) {
          return AccountsService.listTeamMembers($stateParams.teamId).then(null, function(errorData) {
            $state.go('error');
          });
        }],
        users: ['AccountsService', '$state', '$stateParams', function (AccountsService, $state, $stateParams) {
          return AccountsService.list().then(null, function(errorData) {
            $state.go('error');
          });
        }]
      },
      ncyBreadcrumb: {
        skip: true
      }
    })
    .state('dashboard.accounts.users', {
      url: '^/users',
      template: usersTemplate,
      controller: 'UsersController as vm',
      authenticate: true,
      resolve: {
        users: ['AccountsService', '$state', '$stateParams', function (AccountsService, $state, $stateParams) {
          return AccountsService.list().then(null, function(errorData) {
            $state.go('error');
          });
        }]
      },
      ncyBreadcrumb: {
        skip: true
      }
    })
    .state('dashboard.accounts.teamSettings', {
      url: '^/teams/{teamId}/settings',
      template: teamSettingsTemplate,
      controller: 'TeamSettingsController as vm',
      authenticate: true,
      resolve: {
        team: ['AccountsService', '$state', '$stateParams', function (AccountsService, $state, $stateParams) {
          return AccountsService.getTeam($stateParams.teamId).then(null, function(errorData) {
            $state.go('error');
          });
        }]
      },
      ncyBreadcrumb: {
        skip: true
      }
    })
    .state('dashboard.accounts.teamPermissions', {
      url: '^/teams/{teamId}/permissions',
      template: teamPermissionsTemplate,
      controller: 'TeamPermissionsController as vm',
      authenticate: true,
      resolve: {
        team: ['AccountsService', '$state', '$stateParams', function (AccountsService, $state, $stateParams) {
          return AccountsService.getTeam($stateParams.teamId).then(null, function(errorData) {
            $state.go('error');
          });
        }],
        permissions: ['AccountsService', '$state', '$stateParams', function (AccountsService, $state, $stateParams) {
          return AccountsService.getPermissionsForTeam($stateParams.teamId).then(null, function(errorData) {
            $state.go('error');
          });
        }]
      },
      ncyBreadcrumb: {
        skip: true
      }
    });
}

module.exports = getRoutes;
