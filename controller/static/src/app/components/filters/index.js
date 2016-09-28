'use strict';

var angular = require('angular');

module.exports = angular
  .module('ducp.filters', [])
  .filter('split', require('./split'))
  .filter('capitalize', require('./capitalize'))
  .filter('roleDisplay', require('./roledisplay'))
  .filter('fromCalendar', require('./fromcalendar'))
  .filter('isEmpty', function () {
        var bar;
        return function (obj) {
            for (bar in obj) {
                if (obj.hasOwnProperty(bar)) {
                    return false;
                }
            }
            return true;
        };
    });
