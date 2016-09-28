'use strict';

var angular = require('angular');

module.exports = angular.module('ducp.directives', [])
  .directive('fileInput', require('./fileinput'))
  .directive('message', require('./message'))
  .directive('selectableText', require('./selectableText'))
  .directive('stopEvent', require('./stopEvent'))
  .directive('convertToNumber', require('./convertToNumber'))
  .directive('convertToMillis', require('./convertToMillis'))
  .directive('convertToNanos', require('./convertToNanos'))
  .directive('daysToNanos', require('./daysToNanos'))
  .directive('convertToBoolean', require('./convertToBoolean'))
  .directive('uiSrefIf', require('./uiSrefIf'))
  .directive('jquery', require('./jquery'))
  .directive('svgIcon', require('./svgIcon'));

