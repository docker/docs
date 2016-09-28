'use strict';

require('./ng-table.css');

var angular = require('angular');

ConfigureNgTable.$inject = ['$templateCache'];
function ConfigureNgTable($templateCache) {
  // Set up ng-table pager template
  $templateCache.put('custom/pager', require('./pager.html'));
  $templateCache.put('custom/blank', '');
}

module.exports = angular.module('ucpNgTable', []).run(ConfigureNgTable);
