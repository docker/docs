'use strict';

var $ = require('jquery');

MessageDirective.$inject = ['$timeout', 'MessageService'];
function MessageDirective($timeout, MessageService) {
  function dismiss(msg) {
    MessageService.dismissMessage(msg);
  }
  return {
    restrict: 'E',
    replace: true,
    scope: {
      ngModel: '='
    },
    template: '<div class="ui {{ ngModel.severity }} message">' +
                '<i class="close icon" ng-click="dismiss(ngModel)"></i>' +
                '<div class="header" ng-bind-html="ngModel.title"></div>' +
                '<p ng-bind-html="ngModel.message"></p>' +
                '<div class="ui divider" ng-if="ngModel.log"></div>' +
                '<strong ng-if="ngModel.log">Log</strong>' +
                '<pre ng-if="ngModel.log" class="logs">{{ ngModel.log }}</pre>' +
              '</div>',
    link: function(scope, element, attrs) {
      scope.dismiss = dismiss;

      // Only fade out success messages
      if(scope.ngModel.severity === 'success') {
        $timeout(function(){
            element.hide();
          },
          5000
        );
      }
    }
  };
}

module.exports = MessageDirective;
