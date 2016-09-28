'use strict';

CertsService.$inject = ['$http'];
function CertsService($http) {
  return {
    getCerts: function() {
      var promise = $http.get('/api/nodes/certs')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    updateCerts: function(certs) {
      return $http({
        method: 'POST',
        url: '/api/nodes/certs',
        data: certs
      });
    }
  };
}

module.exports = CertsService;
