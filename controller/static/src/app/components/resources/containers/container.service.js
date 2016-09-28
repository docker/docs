'use strict';

ContainerService.$inject = ['$http'];
function ContainerService($http) {
  return {
    list: function(all) {
      var url = '/containers/json';
      if(all) {
        url = url + '?all=1';
      }
      var promise = $http.get(url)
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    inspect: function(containerId) {
      var promise = $http
        .get('/containers/' + containerId + '/json')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    'top': function(containerId) {
      var promise = $http
        .get('/containers/' + containerId + '/top')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    stats: function(containerId) {
      var promise = $http
        .get('/containers/' + containerId + '/stats')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    remove: function(containerId, removeVolumes) {
      var removeVolumeOpt = removeVolumes ? 1 : 0;
      var promise = $http
        .delete('/containers/' + containerId + '?v=' + removeVolumeOpt + '&force=1')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    stop: function(containerId) {
      var promise = $http
        .post('/containers/' + containerId + '/stop')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    restart: function(containerId) {
      var promise = $http
        .post('/containers/' + containerId + '/restart')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    scale: function(containerId, numOfInstances) {
      var promise = $http
        .post('/api/containers/' + containerId + '/scale?n=' + numOfInstances)
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    rename: function(old, newName) {
      var promise = $http
        .post('/containers/' + old + '/rename?name=' + newName)
        .then(function(response) {
          return response.data;
        });
      return promise;
    }
  };
}

module.exports = ContainerService;
