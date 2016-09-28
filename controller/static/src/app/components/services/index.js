'use strict';

var angular = require('angular');

module.exports = angular.module('ducp.core', [])
  .factory('CAService', require('./ca.service'))
  .factory('LoaderService', require('./loader.service'))
  .factory('MessageService', require('./message.service'))
  .factory('AuthService', require('./auth.service'));
