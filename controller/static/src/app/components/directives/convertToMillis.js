'use strict';

function ConvertToMillis() {
  return {
    require: 'ngModel',
    link: function(scope, element, attrs, ngModel) {
      ngModel.$parsers.push(function(val) {
        return val * 1000;
      });
      ngModel.$formatters.push(function(val) {
        return val / 1000;
      });
    }
  };
}

module.exports = ConvertToMillis;
