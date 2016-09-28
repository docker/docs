'use strict';

var settingsTemplate = require('./settingsBase.html');
var loggingTemplate = require('./logging.html');
var licensingTemplate = require('./licensing.html');
var certsTemplate = require('./certs.html');
var swarmTemplate = require('./swarm.html');
var registryTemplate = require('./registry.html');
var trackingTemplate = require('./tracking.html');
var enziTemplate = require('./enzi.html');
var schedulerTemplate = require('./scheduler.html');
var trustTemplate = require('./trust.html');

getRoutes.$inject = ['$stateProvider', '$urlRouterProvider'];
function getRoutes($stateProvider, $urlRouterProvider) {
  $stateProvider
    .state('dashboard.settings', {
      redirectTo: 'dashboard.settings.swarm',
      url: '^/settings',
      template: settingsTemplate,
      authenticate: true,
      ncyBreadcrumb: {
        label: 'Admin Settings'
      }
    })
    .state('dashboard.settings.logging', {
      url: '/logging',
      template: loggingTemplate,
      controller: 'LoggingController',
      controllerAs: 'vm',
      authenticate: true,
      ncyBreadcrumb: {
        skip: true
      },
      resolve: {
        remoteLoggingConfig: ['LoggingService', '$state', function(LoggingService, $state) {
          return LoggingService.getRemoteLoggingConfig().then(null, function(errorData) {
            $state.go('error');
          });
        }]
      }
    })
    .state('dashboard.settings.swarm', {
      url: '/swarm',
      template: swarmTemplate,
      controller: 'SwarmController',
      controllerAs: 'vm',
      authenticate: true,
      ncyBreadcrumb: {
        skip: true
      }
    })
    .state('dashboard.settings.licensing', {
      url: '/licensing',
      template: licensingTemplate,
      controller: 'LicensingController',
      controllerAs: 'vm',
      authenticate: true,
      ncyBreadcrumb: {
        skip: true
      }
    })
    .state('dashboard.settings.certs', {
      url: '/certs',
      template: certsTemplate,
      controller: 'CertsController',
      controllerAs: 'vm',
      authenticate: true,
      ncyBreadcrumb: {
        skip: true
      }
    })
    .state('dashboard.settings.registry', {
      url: '/registry',
      template: registryTemplate,
      controller: 'RegistrySettingsController',
      controllerAs: 'vm',
      authenticate: true,
      ncyBreadcrumb: {
        skip: true
      }
    })
    .state('dashboard.settings.tracking', {
      url: '/usage',
      template: trackingTemplate,
      controller: 'TrackingController',
      controllerAs: 'vm',
      authenticate: true,
      ncyBreadcrumb: {
        skip: true
      },
      resolve: {
        trackingConfig: ['TrackingService', '$state', function(TrackingService, $state) {
          return TrackingService.getTrackingConfig().then(null, function(errorData) {
            $state.go('error');
          });
        }]
      }
    })
    .state('dashboard.settings.scheduler', {
      url: '/scheduler',
      template: schedulerTemplate,
      controller: 'SchedulerController',
      controllerAs: 'vm',
      authenticate: true,
      ncyBreadcrumb: {
        skip: true
      },
      resolve: {
        schedulerConfig: ['SchedulerService', '$state', function(SchedulerService, $state) {
          return SchedulerService.getSchedulerConfig().then(null, function(errorData) {
            $state.go('error');
          });
        }]
      }
    })
    .state('dashboard.settings.trust', {
      url: '/trust',
      template: trustTemplate,
      controller: 'TrustController',
      controllerAs: 'vm',
      authenticate: true,
      ncyBreadcrumb: {
        skip: true
      },
      resolve: {
        trustConfig: ['TrustService', '$state', function(TrustService, $state) {
          return TrustService.getTrustConfig().then(null, function(errorData) {
            $state.go('error');
          });
        }]
      }
    })
    .state('dashboard.settings.enzi', {
      url: '/enzi',
      template: enziTemplate,
      controller: 'EnziController',
      controllerAs: 'vm',
      authenticate: true,
      ncyBreadcrumb: {
        skip: true
      }
    });
}

module.exports = getRoutes;
