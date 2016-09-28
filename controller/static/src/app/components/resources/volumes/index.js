'use strict';

var angular = require('angular');

module.exports = angular.module('ducp.volumes', [])
       .controller('VolumesController', require('./volumes.controller'))
       .factory('VolumesService', require('./volumes.service'));
