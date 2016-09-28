'use strict';

var angular = require('angular');

// TODO: Break this up into a series of directive and filter component files
module.exports = angular.module('ducp.core', [])
  .directive('resetField', ['$compile', function($compile) {
    return {
      require: 'ngModel',
      scope: {
      },
      link: function(scope, element, attrs, ctrl) {
        var template = $compile('<i class="delete link icon"></i>')(scope);
        element.after(template);

        element.parent().find('i').bind('click', function(e) {
          ctrl.$setViewValue('');
          ctrl.$render();
          setTimeout(function() {
            element[0].focus();
          }, 0, false);
          scope.$apply();
        });
      }
    };
  }])
  .filter('cut', function () {
    return function (value, wordwise, max, tail) {
      if (!value) {
        return '';
      }

      max = parseInt(max, 10);
      if (!max) {
        return value;
      }

      if (value.length <= max) {
        return value;
      }

      value = value.substr(0, max);
      if (wordwise) {
        var lastspace = value.lastIndexOf(' ');
        if (lastspace !== -1) {
          value = value.substr(0, lastspace);
        }
      }

      return value + (tail || ' …');
    };
  })
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
