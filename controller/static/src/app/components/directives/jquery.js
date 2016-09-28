'use strict';

// Broadcasts a 'ngRepeatFinished' event used to trigger jquery calls
function JQueryDirective() {
  return {
    link: function(scope, element, attrs) {
      if (scope.$last) {
        setTimeout(function() {
          scope.$emit('ngRepeatFinished', element, attrs);
        }, 0);
      }
    }
  };
}

module.exports = JQueryDirective;
