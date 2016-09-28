'use strict';

NodeService.$inject = ['$http', '$q'];
function NodeService($http, $q) {

  function updateNode(id, version, spec) {
    return $http({
      method: 'POST',
      url: '/nodes/' + id + '/update?version=' + version,
      data: spec
    });
  }

  function emptyResponse(response) {
    return $q.reject({
      data: {
        message: 'Unable to retrieve current node status'
      }
    });
  }

  return {
    remove: function(id) {
      return $http({
        method: 'DELETE',
        url: '/nodes/' + id
      });
    },
    accept: function(id) {
      var promise = $http.get('/nodes/' + id)
        .then(function(response) {
          if(!response.data) {
            return emptyResponse(response);
          }

          var spec = response.data.Spec;
          spec.Membership = 'accepted';
          return updateNode(id, response.data.Version.Index, spec);
        });
      return promise;
    },
    activate: function(id) {
      var promise = $http.get('/nodes/' + id)
        .then(function(response) {
          if(!response.data) {
            return emptyResponse(response);
          }
          var spec = response.data.Spec;
          spec.Availability = 'active';
          return updateNode(id, response.data.Version.Index, spec);
        });
      return promise;
    },
    pause: function(id) {
      var promise = $http.get('/nodes/' + id)
        .then(function(response) {
          if(!response.data) {
            return emptyResponse(response);
          }

          var spec = response.data.Spec;
          spec.Availability = 'pause';
          return updateNode(id, response.data.Version.Index, spec);
        });
      return promise;
    },
    drain: function(id) {
      var promise = $http.get('/nodes/' + id)
        .then(function(response) {
          if(!response.data) {
            return emptyResponse(response);
          }

          var spec = response.data.Spec;
          spec.Availability = 'drain';
          return updateNode(id, response.data.Version.Index, spec);
        });
      return promise;
    },
    promote: function(id) {
      var promise = $http.get('/nodes/' + id)
        .then(function(response) {
          if(!response.data) {
            return emptyResponse(response);
          }

          var spec = response.data.Spec;
          spec.Role = 'manager';
          return updateNode(id, response.data.Version.Index, spec);
        });
      return promise;
    },
    demote: function(id) {
      var promise = $http.get('/nodes/' + id)
        .then(function(response) {
          if(!response.data) {
            return emptyResponse(response);
          }

          var spec = response.data.Spec;
          spec.Role = 'worker';
          return updateNode(id, response.data.Version.Index, spec);
        });
      return promise;
    },
    list: function() {
      var promise = $http.get('/nodes')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    inspect: function(id) {
      var promise = $http.get('/nodes/' + id)
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    listClassic: function() {
      var promise = $http.get('/api/nodes')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    update: updateNode
  };
}

module.exports = NodeService;
