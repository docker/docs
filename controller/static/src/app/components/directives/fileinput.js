'use strict';

function FileInputDirective() {
  return {
    restrict: 'A',
    scope: {
      fileInput: '='
    },
    link: function (scope, element, attributes) {
      element.bind('change', function (changeEvent) {
        var reader = new FileReader();
        reader.onload = function (loadEvent) {
          scope.$apply(function () {
            scope.fileInput = loadEvent.target.result;
          });
        };
        reader.readAsText(changeEvent.target.files[0]);
      });
    }
  };
}

module.exports = FileInputDirective;
