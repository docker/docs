'use strict';

TrustService.$inject = ['$http'];
function TrustService($http) {
  return {
    getTrustConfig: function() {
      var promise = $http.get('/api/config/trust')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    setTrustConfig: function(config) {
      return $http({
        method: 'POST',
        url: '/api/config/trust',
        headers: {
          'Content-Type': 'application/json'
        },
        data: config
      });
    }
  };
}

module.exports = TrustService;
