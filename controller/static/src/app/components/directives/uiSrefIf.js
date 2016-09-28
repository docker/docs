'use strict';

UiSrefIf.$inject = ['$compile'];
function UiSrefIf($compile) {
  return {
    link: function($scope, $element, $attrs) {

      var uiSrefVal = $attrs.uiSrefVal,
        uiSrefIf = $attrs.uiSrefIf;

        $element.removeAttr('ui-sref-if');
        $element.removeAttr('ui-sref-val');

        $scope.$watch(
          function(){
            return $scope.$eval(uiSrefIf);
          },
          function(bool) {
            if (bool) {

              $element.attr('ui-sref', uiSrefVal);
            } else {

              $element.removeAttr('ui-sref');
              $element.removeAttr('href');
            }
            $compile($element)($scope);
          }
        );
    }
  };
}

module.exports = UiSrefIf;
