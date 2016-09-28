'use strict';

LoaderInterceptor.$inject = ['$q', '$rootScope'];
function LoaderInterceptor($q, $rootScope) {
  var service = this;
  var nLoadings = 0;
  service.request = function(request) {
    // If 'noLoader' is set on the http request, then don't trigger a spinner
    if(request && request.noLoader) {
      return request;
    }

    nLoadings += 1;
    $rootScope.isLoadingView = true;

    return request;
  };
  service.response = function(response) {
    // If 'noLoader' is set on the http request, then don't trigger a spinner
    if(response.config && response.config.noLoader) {
      return response;
    }

    nLoadings -= 1;
    if (nLoadings === 0) {
      $rootScope.isLoadingView = false;
    }

    return response;
  };
  service.responseError = function(response) {
    // If 'noLoader' is set on the http request, then don't trigger a spinner
    if(response.config && response.config.noLoader) {
      return $q.reject(response);
    }

    nLoadings -= 1;
    if (!nLoadings) {
      $rootScope.isLoadingView = false;
    }

    return $q.reject(response);
  };
}

module.exports = LoaderInterceptor;

