'use strict';

function DaysToNanos() {
  return {
    require: 'ngModel',
    link: function(scope, element, attrs, ngModel) {
      ngModel.$parsers.push(function(val) {
        return val * 1000000000 * 60 * 60 * 24;
      });
      ngModel.$formatters.push(function(val) {
        return val / (1000000000 * 60 * 60 * 24);
      });
    }
  };
}

module.exports = DaysToNanos;
