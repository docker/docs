'use strict';

function ConvertToNanos() {
  return {
    require: 'ngModel',
    link: function(scope, element, attrs, ngModel) {
      ngModel.$parsers.push(function(val) {
        return val * 1000000000;
      });
      ngModel.$formatters.push(function(val) {
        return val / 1000000000;
      });
    }
  };
}

module.exports = ConvertToNanos;
