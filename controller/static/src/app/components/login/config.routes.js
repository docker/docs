'use strict';

var loginTemplate = require('./login.html');
var licenseReminderTemplate = require('./licenseReminder.html');

getRoutes.$inject = ['$stateProvider', '$urlRouterProvider'];
function getRoutes($stateProvider, $urlRouterProvider) {
  $stateProvider
    .state('login', {
      url: '/login?return&returnParams',
      template: loginTemplate,
      controller: 'LoginController',
      controllerAs: 'vm',
      authenticate: false
    })
    .state('licenseReminder', {
      template: licenseReminderTemplate,
      controller: 'LicenseReminderController',
      controllerAs: 'vm',
      authenticate: true
    })
    .state('logout', {
      url: '/logout',
      controller: 'LogoutController',
      controllerAs: 'vm',
      authenticate: true
    });
}

module.exports = getRoutes;
