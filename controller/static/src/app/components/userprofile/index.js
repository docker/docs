'use strict';

var angular = require('angular');

var userController = require('./user.controller.js');
var routes = require('./routes.js');

module.exports = angular.module('ducp.user', ['ducp.accounts'])
  .controller('UserController', userController)
  .config(routes);
