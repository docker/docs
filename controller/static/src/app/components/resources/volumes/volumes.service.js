'use strict';

VolumesService.$inject = ['$http'];
function VolumesService($http) {
  return {
    list: function() {
      var promise = $http
      .get('/volumes')
      .then(function(response) {
        return response.data.Volumes;
      });
      return promise;
    },
    create: function(request) {
      var promise = $http
      .post('/volumes/create', request)
      .then(function(response) {
        return response.data;
      });
      return promise;
    },
    remove: function(volumeName) {
      var promise = $http
      .delete('/volumes/' + volumeName)
      .then(function(response) {
        return response.data;
      });
      return promise;
    }
  };
}

module.exports = VolumesService;
