'use strict';

NetworksService.$inject = ['$http'];
function NetworksService($http) {
  return {
    list: function() {
      var promise = $http
      .get('/networks')
      .then(function(response) {
        return response.data;
      });
      return promise;
    },
    create: function(data) {
      var promise = $http
      .post('/networks/create', data)
      .then(function(response) {
        return response.data;
      });
      return promise;
    },
    remove: function(network) {
      var promise = $http
      .delete('/networks/' + network.Id)
      .then(function(response) {
        return response.data;
      });
      return promise;
    }
  };
}

module.exports = NetworksService;
