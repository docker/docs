'use strict';

var errorTemplate = require('./components/error/error.html');

getRoutes.$inject = ['$stateProvider', '$urlRouterProvider'];

function getRoutes($stateProvider, $urlRouterProvider) {
  $urlRouterProvider.otherwise(function ($injector) {
    var $state = $injector.get('$state');
    $state.go('dashboard.main');
  });
}

module.exports = getRoutes;
