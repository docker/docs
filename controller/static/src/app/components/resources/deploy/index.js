'use strict';

var angular = require('angular');

var routes = require('./config.routes');
var deployWizardController = require('./deploy.wizard.controller');
var deployYmlController = require('./deploy.yml.controller');
var deployStackController = require('./deploy.stack.controller');

require('./deploy.css');

module.exports = angular.module('ducp.deploy', [])
  .controller('DeployWizardController', deployWizardController)
  .controller('DeployYmlController', deployYmlController)
  .controller('DeployStackController', deployStackController)
  .config(routes);
