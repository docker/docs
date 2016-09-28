'use strict';

var angular = require('angular');
var footerController = require('./footer.controller');
var routes = require('./config.routes');
require('./layout.css');

module.exports = angular.module('ducp.layout', [])
  .controller('FooterController', footerController)
  .config(routes);
