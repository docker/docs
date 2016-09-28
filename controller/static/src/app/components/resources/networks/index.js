'use strict';

var angular = require('angular');

var networksService = require('./networks.service');
var networksController = require('./networks.controller');

require('./networks.css');

module.exports = angular.module('ducp.networks', [])
  .factory('NetworksService', networksService)
  .controller('NetworksController', networksController);
