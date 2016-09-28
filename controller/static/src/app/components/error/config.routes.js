'use strict';

var errorTemplate = require('./error.html');

getRoutes.$inject = ['$stateProvider', '$urlRouterProvider'];

function getRoutes($stateProvider, $urlRouterProvider) {
  $stateProvider
    .state('error', {
      template: errorTemplate,
      authenticate: false
    });
}

module.exports = getRoutes;
