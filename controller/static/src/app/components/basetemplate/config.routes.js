'use strict';

var template = require('./base.html');

getRoutes.$inject = ['$stateProvider', '$urlRouterProvider'];
function getRoutes($stateProvider, $urlRouterProvider) {
  $stateProvider
    .state('dashboard', {
      url: '/',
      abstract: true,
      template: template,
      controller: 'BaseController as vm'
    });
}


module.exports = getRoutes;
