'use strict';

var angular = require('angular');
require('./images.css');

module.exports = angular.module('ducp.images', [])
  .factory('ImagesService', require('./images.service'))
  .controller('ImagesController', require('./images.controller'));
