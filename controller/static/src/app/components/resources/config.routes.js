'use strict';

var template = require('./resources.html');

var applicationsTemplate = require('./applications/applications.html');
var containersTemplate = require('./containers/containers.html');
var imagesTemplate = require('./images/images.html');
var networksTemplate = require('./networks/networks.html');
var volumesTemplate = require('./volumes/volumes.html');

var serviceTemplate = require('./services/services.html');
var serviceInspectTemplate = require('./services/inspect.service.html');
var nodesTemplate = require('./nodes/nodes.html');
var nodeInspectTemplate = require('./nodes/inspect.node.html');

getRoutes.$inject = ['$stateProvider', '$urlRouterProvider'];
function getRoutes($stateProvider, $urlRouterProvider) {
  $stateProvider
  .state('dashboard.resources', {
    url: '^/resources',
    redirectTo: 'dashboard.resources.services',
    template: template,
    controller: 'ResourcesController',
    controllerAs: 'vm',
    authenticate: true,
    ncyBreadcrumb: {
      label: 'Resources'
    }
  })
  .state('dashboard.resources.applications', {
    url: '^/resources/applications',
    template: applicationsTemplate,
    controller: 'ApplicationsController',
    controllerAs: 'vm',
    authenticate: true,
    ncyBreadcrumb: {
      label: 'Applications'
    }
  })
  .state('dashboard.resources.services', {
    url: '/services',
    redirectTo: 'dashboard.resources.services.list'
  })
  .state('dashboard.resources.services.list', {
    url: '',
    views: {
      '@dashboard.resources': {
        template: serviceTemplate,
        controller: 'ServicesController',
        controllerAs: 'vm'
      }
    },
    authenticate: true,
    ncyBreadcrumb: {
      parent: 'dashboard.resources',
      label: 'Services'
    },
    resolve: {
      nodes: ['NodeService', '$state', function(NodeService, $state) {
        return NodeService.list().then(null, function(error) {
          $state.go('error');
        });
      }]
    }
  })
  .state('dashboard.resources.services.inspect', {
    url: '^/resources/services/{id}',
    views: {
      '@dashboard.resources': {
        template: serviceInspectTemplate,
        controller: 'InspectServiceController',
        controllerAs: 'vm'
      }
    },
    authenticate: true,
    ncyBreadcrumb: {
      parent: 'dashboard.resources.services.list',
      label: 'Inspect Service'
    },
    resolve: {
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
      nodes: ['NodeService', '$state', function(NodeService, $state) {
        return NodeService.list().then(null, function(error) {
          $state.go('error');
        });
      }]
    }
  })
  .state('dashboard.resources.containers', {
    url: '^/resources/containers',
    template: containersTemplate,
    controller: 'ContainersController',
    controllerAs: 'vm',
    authenticate: true,
    ncyBreadcrumb: {
      label: 'Containers'
    }
  })
  .state('dashboard.resources.images', {
    url: '^/resources/images',
    template: imagesTemplate,
    controller: 'ImagesController',
    controllerAs: 'vm',
    authenticate: true,
    ncyBreadcrumb: {
      label: 'Images'
    }
  })
  .state('dashboard.resources.nodes', {
    url: '/nodes',
    redirectTo: 'dashboard.resources.nodes.list'
  })
  .state('dashboard.resources.nodes.list', {
    url: '',
    views: {
      '@dashboard.resources': {
        template: nodesTemplate,
        controller: 'NodesController',
        controllerAs: 'vm'
      }
    },
    authenticate: true,
    ncyBreadcrumb: {
      parent: 'dashboard.resources',
      label: 'Nodes'
    },
    resolve: {
      info: ['DashboardService', function(DashboardService) {
        return DashboardService.info();
      }]
    }
  })
  .state('dashboard.resources.nodes.inspect', {
    url: '^/resources/nodes/{id}',
    views: {
      '@dashboard.resources': {
        template: nodeInspectTemplate,
        controller: 'InspectNodeController',
        controllerAs: 'vm'
      }
    },
    authenticate: true,
    ncyBreadcrumb: {
      parent: 'dashboard.resources.nodes.list',
      label: 'Inspect Node'
    },
    resolve: {}
  })
  .state('dashboard.resources.networks', {
    url: '^/resources/networks',
    template: networksTemplate,
    controller: 'NetworksController',
    controllerAs: 'vm',
    authenticate: true,
    ncyBreadcrumb: {
      label: 'Networks'
    }
  })
  .state('dashboard.resources.volumes', {
    url: '^/resources/volumes',
    template: volumesTemplate,
    controller: 'VolumesController',
    controllerAs: 'vm',
    authenticate: true,
    ncyBreadcrumb: {
      label: 'Volumes'
    }
  })
  ;
}


module.exports = getRoutes;
