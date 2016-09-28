'use strict';

function ConvertToBoolean() {
  return {
    require: 'ngModel',
    link: function(scope, element, attrs, ngModel) {
      ngModel.$parsers.push(function(val) {
        return val === 'true';
      });
      ngModel.$formatters.push(function(val) {
        return val ? 'true' : 'false';
      });
    }
  };
}

module.exports = ConvertToBoolean;
