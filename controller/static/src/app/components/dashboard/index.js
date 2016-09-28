'use strict';

var angular = require('angular');

var serviceService = require('../resources/services/service.service');
var taskService = require('../resources/services/task.service');

require('./dashboard.css');

module.exports = angular
  .module('ducp.dashboard', [])
  .factory('DashboardService', require('./dashboard.service.js'))
  .factory('ServiceService', serviceService)
  .factory('TaskService', taskService)
  .controller('DashboardController', require('./dashboard.controller.js'))
  .config(require('./config.routes.js'));
