'use strict';

var angular = require('angular');

var routes = require('./config.routes');

require('./licenseReminder.css');
require('./login.css');

module.exports = angular.module('ducp.login', [])
  .controller('LogoutController', require('./logout.controller'))
  .controller('LoginController', require('./login.controller'))
  .controller('LicenseReminderController', require('./licenseReminder.controller'))
  .config(routes);
