'use strict';

var dashboardTemplate = require('./dashboard.html');

getRoutes.$inject = ['$stateProvider', '$urlRouterProvider'];
function getRoutes($stateProvider, $urlRouterProvider) {
  $stateProvider
    .state('dashboard.main', {
      url: '^/dashboard',
      template: dashboardTemplate,
      controller: 'DashboardController',
      controllerAs: 'vm',
      authenticate: true,
      ncyBreadcrumb: {
        label: 'Dashboard'
      }
    });
}

module.exports = getRoutes;
