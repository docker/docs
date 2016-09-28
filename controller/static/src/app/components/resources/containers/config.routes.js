'use strict';

var containersTemplate = require('./containers.html');
var execTemplate = require('./exec.html');
var statsTemplate = require('./stats.html');
var logsTemplate = require('./logs.html');
var inspectBaseTemplate = require('./inspectBase.html');
var inspectDetailsTemplate = require('./inspect.html');

require('./inspect.css');
require('./logs.css');

getRoutes.$inject = ['$stateProvider', '$urlRouterProvider'];
function getRoutes($stateProvider, $urlRouterProvider) {
  $stateProvider
  .state('dashboard.inspect', {
    url: '^/resources/containers/{id}',
    redirectTo: 'dashboard.inspect.details',
    ncyBreadcrumb: {
      parent: 'dashboard.resources.containers',
      label: 'Container'
    },
    template: inspectBaseTemplate,
    controller: 'InspectBaseController as vm',
    authenticate: true,
    resolve: {
      resolvedContainer: ['ContainerService', '$state', '$stateParams', function(ContainerService, $state, $stateParams) {
        return ContainerService.inspect($stateParams.id).then(null, function(errorData) {
          $state.go('error');
        });
      }]
    }
  })
  .state('dashboard.inspect.details', {
    url: '/inspect',
    ncyBreadcrumb: {
      skip: true
    },
    authenticate: true,
    template: inspectDetailsTemplate,
    controller: 'InspectDetailsController as vm',
    resolve: {
      resolvedContainer: ['ContainerService', '$state', '$stateParams', function(ContainerService, $state, $stateParams) {
        return ContainerService.inspect($stateParams.id).then(null, function(errorData) {
          $state.go('error');
        });
      }],
      fromState: ['$rootScope', function($rootScope) {
        return $rootScope.toState;
      }],
      fromParams: ['$rootScope', function($rootScope) {
        return $rootScope.toParams;
      }]
    }
  })
  .state('dashboard.inspect.logs', {
    url: '/logs',
    ncyBreadcrumb: {
      skip: true
    },
    template: logsTemplate,
    controller: 'LogsController',
    controllerAs: 'vm',
    authenticate: true,
    resolve: {
      resolvedContainer: ['ContainerService', '$state', '$stateParams', function(ContainerService, $state, $stateParams) {
        return ContainerService.inspect($stateParams.id).then(null, function(errorData) {
          $state.go('error');
        });
      }]
    }
  })
  .state('dashboard.inspect.stats', {
    url: '/stats',
    ncyBreadcrumb: {
      skip: true
    },
    template: statsTemplate,
    controller: 'ContainerStatsController',
    controllerAs: 'vm',
    authenticate: true
  })
  .state('dashboard.inspect.exec', {
    url: '/exec',
    ncyBreadcrumb: {
      skip: true
    },
    template: execTemplate,
    controller: 'ExecController',
    controllerAs: 'vm',
    authenticate: true
  });
}

module.exports = getRoutes;
