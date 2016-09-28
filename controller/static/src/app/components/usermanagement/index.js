'use strict';

var angular = require('angular');

module.exports = angular.module('ducp.accounts', [])
  .factory('AccountsService', require('./accounts.service'))
  .controller('AccountsBaseController', require('./accountsBase.controller'))
  .controller('UsersController', require('./users.controller'))
  .controller('TeamController', require('./team.controller'))
  .controller('TeamPermissionsController', require('./teamPermissions.controller'))
  .controller('TeamSettingsController', require('./teamSettings.controller'))
  .config(require('./config.routes'));
