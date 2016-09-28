'use strict';

SchedulerService.$inject = ['$http'];
function SchedulerService($http) {
  return {
    getSchedulerConfig: function() {
      var promise = $http.get('/api/config/scheduling')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    setSchedulerConfig: function(config) {
      return $http({
        method: 'POST',
        url: '/api/config/scheduling',
        headers: {
          'Content-Type': 'application/json'
        },
        data: config
      });
    }
  };
}

module.exports = SchedulerService;
