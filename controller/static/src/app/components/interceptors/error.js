'use strict';

ErrorInterceptor.$inject = ['$q', '$rootScope', 'MessageService'];
function ErrorInterceptor($q, $rootScope, MessageService) {
  var service = this;
  service.responseError = function(response) {
    if(response.status === 401 && !response.config.ignore401) {
      $rootScope.$state.go('login');
    } else if(response.status >= 400) {
      return $q.reject(response);
    }
    return response;
  };
}

module.exports = ErrorInterceptor;
