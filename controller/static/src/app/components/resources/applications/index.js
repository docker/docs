'use strict';

var angular = require('angular');

module.exports = angular.module('ducp.applications', [])
  .controller('ApplicationsController', require('./applications.controller.js'))
  .factory('ApplicationsService', require('./service.js'));
