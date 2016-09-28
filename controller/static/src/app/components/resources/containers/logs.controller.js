'use strict';

var _ = require('lodash');

LogsController.$inject = ['resolvedContainer', '$http', '$stateParams', 'LogService', 'LogConstants', 'MessageService', '$timeout', '$scope', '$state'];
function LogsController(resolvedContainer, $http, $stateParams, LogService, LogConstants, MessageService, $timeout, $scope, $state) {
  $state.current.ncyBreadcrumb.parent = 'dashboard.inspect.details({id: \'' + $stateParams.id + '\'})';
  $scope.LogConstants = LogConstants;

  var vm = this;
  vm.id = $stateParams.id;
  vm.logs = '';
  vm.logParams = LogService.defaultParams;
  vm.toggleTail = toggleTail;

  var url = window.location.href;
  var urlparts = url.split('/');
  var scheme = urlparts[0];
  var wsScheme = 'ws';
  var websocket;

  if (scheme === 'https:') {
    wsScheme = 'wss';
  }

  function toggleTail() {
    if(vm.logParams.tail === LogConstants.DISABLE_TAIL) {
      vm.logParams.tail = LogConstants.TAIL_DEFAULT;
    } else {
      vm.logParams.tail = LogConstants.DISABLE_TAIL;
    }
  }

  function closeWebsocket() {
    if(websocket) {
      websocket.close();
    }
    vm.logs = '';
  }

  function getLogs(params) {
    closeWebsocket();
    LogService.get(vm.id, params)
      .then(function(ws) {
        websocket = ws;
        websocket.onopen = function(evt) {
          websocket.onmessage = function(msg) {
            vm.logs += msg.data;
            $timeout(function() {
              $scope.$apply();
            });
          };
        };
      }, function(error) {
        MessageService.addErrorMessage('Could not retrieve logs');
      });
  }

  $scope.$watch('vm.logParams', function(after, before) {
    // If the config has actually changed. re-request logs
    if(!_.isEqual(after, before)) {
      getLogs(after);
    }
  }, true);

  $scope.$on('$destroy', function() {
    closeWebsocket();
  });

  // Fetch logs
  getLogs(vm.logParams);
}

module.exports = LogsController;
