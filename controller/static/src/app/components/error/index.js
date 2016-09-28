'use strict';

var angular = require('angular');

module.exports = angular
  .module('ducp.error', [])
  .config(require('./config.routes.js'));

