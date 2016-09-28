'use strict';

var angular = require('angular');

var containerService = require('./container.service');
var inspectBaseController = require('./inspectBase.controller');
var inspectDetailsController = require('./inspect.controller');
var containersController = require('./containers.controller');
var containerStatsController = require('./stats.controller');
var containerExecController = require('./exec.controller');
var containerLogsController = require('./logs.controller');
var logService = require('./log.service');
var routes = require('./config.routes');

require('./containers.css');
require('./exec.css');

module.exports = angular.module('ducp.containers', [])
  .factory('ContainerService', containerService)
  .constant('LogConstants', {
    DISABLE_TAIL: 'all',
    TAIL_DEFAULT: 100
  })
  .factory('LogService', logService)
  .controller('InspectBaseController', inspectBaseController)
  .controller('InspectDetailsController', inspectDetailsController)
  .controller('ContainersController', containersController)
  .controller('ContainerStatsController', containerStatsController)
  .controller('ExecController', containerExecController)
  .controller('LogsController', containerLogsController)
  .config(routes);
