'use strict';

ApplicationsService.$inject = ['$http'];
function ApplicationsService($http) {
  return {
    list: function() {
      var promise = $http
      .get('/api/applications')
      .then(function(response) {
        return response.data;
      });
      return promise;
    },
    get: function(name) {
      var promise = $http
        .get('/api/applications/' + name)
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    containers: function(app) {
      var promise = $http
        .get('/api/applications/' + app.name + '/containers')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    remove: function(app) {
      var promise = $http
        .delete('/api/applications/' + app.name)
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    restart: function(app) {
      var promise = $http
        .post('/api/applications/' + app.name + '/restart')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    stop: function(app) {
      var promise = $http
        .post('/api/applications/' + app.name + '/stop')
        .then(function(response) {
          return response.data;
        });
      return promise;
    }
  };
}

module.exports = ApplicationsService;
