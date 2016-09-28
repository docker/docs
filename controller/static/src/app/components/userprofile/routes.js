'use strict';

var userTemplate = require('./user.html');

getRoutes.$inject = ['$stateProvider', '$urlRouterProvider'];
function getRoutes($stateProvider, $urlRouterProvider) {
  $stateProvider
    .state('dashboard.user', {
      url: '^/user',
      template: userTemplate,
      controller: 'UserController',
      controllerAs: 'vm',
      authenticate: 'true',
      ncyBreadcrumb: {
        label: 'User Settings'
      },
      resolve: {
        account: ['AccountsService', '$state', 'AuthService', function (AccountsService, $state, AuthService) {
          return AccountsService.getAccount(AuthService.getCurrentUser().username).then(null, function(errorData) {
            $state.go('error');
          });
        }],
        permissions: ['AccountsService', '$state', 'AuthService', function (AccountsService, $state, AuthService) {
          return AccountsService.getPermissionsForAccount(AuthService.getCurrentUser().username).then(null, function(errorData) {
            $state.go('error');
          });
        }]
      }
    });
}

module.exports = getRoutes;
