'use strict';

var angular = require('angular');

require('./settings.css');

module.exports = angular.module('ducp.settings', [])
  .controller('LoggingController', require('./logging.controller'))
  .factory('LoggingService', require('./logging.service'))
  .controller('TrackingController', require('./tracking.controller'))
  .factory('TrackingService', require('./tracking.service'))
  .controller('SchedulerController', require('./scheduler.controller'))
  .factory('SchedulerService', require('./scheduler.service'))
  .controller('TrustController', require('./trust.controller'))
  .factory('TrustService', require('./trust.service'))
  .factory('EnziService', require('./enzi.service'))
  .controller('EnziController', require('./enzi.controller'))
  .controller('LicensingController', require('./licensing.controller'))
  .factory('LicensingService', require('./licensing.service'))
  .controller('CertsController', require('./certs.controller'))
  .factory('CertsService', require('./certs.service'))
  .controller('SwarmController', require('./swarm.controller'))
  .factory('SwarmService', require('./swarm.service'))
  .controller('RegistrySettingsController', require('./registry.controller'))
  .factory('RegistryService', require('./registry.service'))
  .config(require('./config.routes'));
