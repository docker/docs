'use strict';

// This is designed to prevent ng-clicks from being triggered if the user is selecting
// text by dragging or double clicking

SelectableText.$inject = ['$window', '$timeout'];
function SelectableText($window, $timeout) {
  var i = 0;
  return {
    restrict: 'A',
    priority: 1,
    compile: function (tElem, tAttrs) {
      var fn = '$$clickOnNoSelect' + i++;
      var _ngClick = tAttrs.ngClick;

      tAttrs.ngClick = fn + '($event)';

      return function(scope) {
        var lastAnchorOffset, lastFocusOffset, timer;

        scope[fn] = function(event) {
          var selection = $window.getSelection();
          var anchorOffset = selection.anchorOffset;
          var focusOffset = selection.focusOffset;

          if(focusOffset - anchorOffset !== 0) {
            if(!(lastAnchorOffset === anchorOffset && lastFocusOffset === focusOffset)) {
              lastAnchorOffset = anchorOffset;
              lastFocusOffset = focusOffset;
              if(timer) {
                $timeout.cancel(timer);
                timer = null;
              }
              return;
            }
          }
          lastAnchorOffset = null;
          lastFocusOffset = null;

          // delay invoking click so as to watch for user double-clicking
          // to select words
          timer = $timeout(function() {
            scope.$eval(_ngClick, {$event: event});
            timer = null;
          }, 250);
        };
      };
    }
  };
}

module.exports = SelectableText;
