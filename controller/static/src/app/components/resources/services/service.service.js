'use strict';

ServiceService.$inject = ['$http'];
function ServiceService($http) {
  return {
    list: function() {
      var promise = $http({url: '/services'})
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    inspect: function(id) {
      var promise = $http.get('/services/' + id)
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    create: function(service) {
      var promise = $http.post('/services/create', service)
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    remove: function(id) {
      var promise = $http.delete('/services/' + id)
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    update: function(id, serviceSpec, version) {
      return $http.post('/services/' + id + '/update?version=' + version, serviceSpec)
        .then(function(updateResponse) {
          return updateResponse.data;
        });
    },
    scale: function(id, count) {
      var promise = $http.get('/services/' + id)
        .then(function(response) {
          var service = response.data.Spec;
          service.Mode.Replicated.Replicas = count;

          return $http.post('/services/' + id + '/update?version=' + response.data.Version.Index, service)
            .then(function(updateResponse) {
              return updateResponse.data;
            });
        });
      return promise;
    }
  };
}

module.exports = ServiceService;
