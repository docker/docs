'use strict';

// This directive is designed to stop events from propagating up to parent elements.
// e.g. if a button is clicked, prevent the click from being propagated to parent
// elements after, since this can cause unwanted behaviours.
//
StopEvent.$inject = [];
function StopEvent() {
  return {
    restrict: 'A',
    link: function (scope, element, attr) {
      if(attr && attr.stopEvent) {
        element.bind(attr.stopEvent, function (e) {
          e.stopPropagation();
        });
      }
    }
  };
}

module.exports = StopEvent;
