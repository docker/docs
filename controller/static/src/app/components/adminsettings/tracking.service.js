'use strict';

TrackingService.$inject = ['$http'];
function TrackingService($http) {
  return {
    getTrackingConfig: function() {
      var promise = $http.get('/api/config/tracking')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    setTrackingConfig: function(config) {
      return $http({
        method: 'POST',
        url: '/api/config/tracking',
        headers: {
          'Content-Type': 'application/json'
        },
        data: config
      });
    }
  };
}

module.exports = TrackingService;
