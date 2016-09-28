'use strict';

var angular = require('angular');

require('./resources.css');

module.exports = angular
  .module('ducp.resources', [
    require('./deploy').name,
    require('./applications').name,
    require('./services').name,
    require('./containers').name,
    require('./images').name,
    require('./volumes').name,
    require('./nodes').name,
    require('./networks').name
  ])
  .controller('ResourcesController', require('./resources.controller'))
  .config(require('./config.routes'));
