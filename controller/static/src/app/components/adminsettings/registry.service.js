'use strict';

RegistryService.$inject = ['$http'];
function RegistryService($http) {
  return {
    getRegistry: function() {
      var promise = $http.get('/api/config/registry')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    updateRegistry: function(registry) {
      var promise = $http({
        method: 'POST',
        url: '/api/config/registry',
        data: registry
      });
      return promise;
    }
  };
}

module.exports = RegistryService;
