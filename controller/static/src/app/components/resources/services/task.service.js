'use strict';

var _ = require('lodash');

TaskService.$inject = ['$http'];
function TaskService($http) {
  var stateSummaryMap = {
    new: 'inactive',
    allocated: 'inactive',
    complete: 'inactive',
    shutdown: 'inactive',
    remove: 'inactive',
    dead: 'inactive',
    pending: 'inactive',
    assigned: 'inactive',
    ready: 'inactive',
    accepted: 'inactive',
    preparing: 'updating',
    starting: 'updating',
    running: 'active',
    failed: 'errored',
    rejected: 'errored'
  };

  function listTasks(all) {
    var promise = $http({
      method: 'GET',
      url: '/tasks',
      noLoader: true
    })
    .then(function(response) {
      return _.map(response.data, function(t) {
        if (stateSummaryMap[t.Status.State] === 'active' || stateSummaryMap[t.Status.State] === 'updating') {
          t._Status = 'RUNNING';
        } else {
          t._Status = 'HISTORICAL';
        }

        return t;
      });
    });

    return promise;
  }

  return {
    list: listTasks,
    listForService: function(serviceId) {
      var promise = listTasks()
        .then(function(response) {
          return _.filter(response.data, function(t) {
            return t.ServiceID === serviceId;
          });
        });
      return promise;
    },
    listForNode: function(nodeId) {
      var promise = listTasks()
        .then(function(response) {
          return _.filter(response, function(t) {
            return t.NodeID === nodeId;
          });
        });
      return promise;
    },
    stateSummaryMap: stateSummaryMap
  };
}

module.exports = TaskService;
