'use strict';

var angular = require('angular');

var serviceService = require('./service.service');
var taskService = require('./task.service');
var servicesController = require('./services.controller');
var inspectServiceController = require('./inspect.service.controller');

require('./inspect.service.css');

module.exports = angular.module('ducp.services', [])
  .factory('ServiceService', serviceService)
  .factory('TaskService', taskService)
  .controller('ServicesController', servicesController)
  .controller('InspectServiceController', inspectServiceController);
