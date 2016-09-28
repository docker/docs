'use strict';

LoggingService.$inject = ['$http'];
function LoggingService($http) {
  return {
    getRemoteLoggingConfig: function() {
      var promise = $http.get('/api/config/logging')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    setRemoteLoggingConfig: function(config) {
      return $http({
        method: 'POST',
        url: '/api/config/logging',
        headers: {
          'Content-Type': 'application/json'
        },
        data: config
      });
    },
    disableRemoteLogging: function() {
      var promise = $http.post('/api/config/logging', '{"host":""}');
      return promise;
    }
  };
}

module.exports = LoggingService;
