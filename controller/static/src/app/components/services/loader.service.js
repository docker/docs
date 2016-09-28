'use strict';

LoaderService.$inject = ['$rootScope', '$timeout'];
function LoaderService($rootScope, $timeout) {
  return {
    setMessage: function(loadingMessage) {
      $timeout(function() {
        $rootScope.loadingMessage = loadingMessage;
      });
    },
    clear: function() {
      $timeout(function() {
        $rootScope.loadingMessage = '';
      });
    },
    start: function() {
      $timeout(function() {
        $rootScope.isLoadingView = true;
      });
    },
    stop: function() {
      $timeout(function() {
        $rootScope.loadingMessage = '';
        $rootScope.isLoadingView = false;
      });
    }
  };
}

module.exports = LoaderService;
