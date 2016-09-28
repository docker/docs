'use strict';

var _ = require('lodash');
var $ = require('jquery');

InspectServiceController.$inject = ['networks', 'nodes', 'TaskService', 'ServiceService', 'MessageService', 'NgTableParams', '$scope', '$state', '$stateParams', '$location', '$timeout', '$interval'];

function InspectServiceController(networks, nodes, TaskService, ServiceService, MessageService, NgTableParams, $scope, $state, $stateParams, $location, $timeout, $interval) {
  var vm = this;

  $scope.Math = Math;
  $scope.$location = $location;
  vm.stateSummaryMap = TaskService.stateSummaryMap;

  vm.networks = _.keyBy(networks, 'Id');
  vm.nodes = _.keyBy(nodes, 'ID');
  vm.numNodes = nodes.length;
  vm.runningTasksVisible = true;
  vm.historicalTasksVisible = false;
  vm.tableParams = new NgTableParams({
    filter: {
      '_Status': 'RUNNING'
    },
    sorting: {
      'Status.State': 'desc',
      Instance: 'asc'
    }
  });

  // Methods
  vm.removeService = removeService;
  vm.showRemoveServiceModal = showRemoveServiceModal;
  vm.scaleService = scaleService;
  vm.showScaleServiceModal = showScaleServiceModal;
  vm.removeEnvVar = removeEnvVar;
  vm.removeLabel = removeLabel;
  vm.removeConstraint = removeConstraint;
  vm.startEdit = startEdit;
  vm.cancelEdit = cancelEdit;
  vm.completeEdit = completeEdit;
  vm.inspectContainer = inspectContainer;
  vm.toggleRunningTasks = toggleRunningTasks;
  vm.toggleHistoricalTasks = toggleHistoricalTasks;

  // Load data
  load();
  var loadTasksPromise = $interval(loadTasks, 2000);
  $scope.$on('$destroy', function() {
    $interval.cancel(loadTasksPromise);
  });


  // ---------------------------

  function inspectContainer(id) {
    if (!id) {
      return;
    }

    $state.go('dashboard.inspect.details', {
      id: id
    });
  }

  function toggleRunningTasks() {
    var filters = vm.tableParams.filter();
    if (filters._Status === 'RUNNING') {
      filters._Status = 'NONE';
      vm.runningTasksVisible = false;
    } else if (filters._Status === 'HISTORICAL') {
      delete filters._Status;
      vm.historicalTasksVisible = true;
      vm.runningTasksVisible = true;
    } else {
      filters._Status = 'HISTORICAL';
      vm.historicalTasksVisible = true;
      vm.runningTasksVisible = false;
    }
    vm.tableParams.filter(filters);
    vm.tableParams.reload();
    vm.tableParams.page(1);
  }

  function toggleHistoricalTasks() {
    var filters = vm.tableParams.filter();
    if (filters._Status === 'HISTORICAL') {
      filters._Status = 'NONE';
      vm.historicalTasksVisible = false;
    } else if (filters._Status === 'RUNNING') {
      delete filters._Status;
      vm.historicalTasksVisible = true;
      vm.runningTasksVisible = true;
    } else {
      filters._Status = 'RUNNING';
      vm.historicalTasksVisible = false;
      vm.runningTasksVisible = true;
    }
    vm.tableParams.filter(filters);
    vm.tableParams.reload();
    vm.tableParams.page(1);
  }

  function removeLabel(label) {
    var index = vm.editLabels.indexOf(label);
    vm.editLabels.splice(index, 1);
  }

  function removeConstraint(c) {
    var index = vm.editConstraints.indexOf(c);
    vm.editConstraints.splice(index, 1);
  }

  function removeEnvVar(envVar) {
    var index = vm.editEnvVars.indexOf(envVar);
    vm.editEnvVars.splice(index, 1);
  }

  function startEdit() {
    vm.editUpdateConfigDelay = null;
    if (vm.service.Spec.UpdateConfig && vm.service.Spec.UpdateConfig.Delay) {
      vm.editUpdateConfigDelay = vm.service.Spec.UpdateConfig.Delay / 1000000000;
    }

    vm.editStopGracePeriod = vm.service.Spec.TaskTemplate.ContainerSpec.StopGracePeriod ? (vm.service.Spec.TaskTemplate.ContainerSpec.StopGracePeriod / 1000000000) : null;

    vm.editEnvVars = _.map(vm.service.Spec.TaskTemplate.ContainerSpec.Env, function(v) {
      return v.split('=');
    });

    vm.editLabels = _.map(vm.service.Spec.Labels, function(value, key) { return { key: key, value: value }; });
    vm.editConstraints = vm.service.Spec.TaskTemplate.Placement.Constraints;

    // Make a copy of the existing spec
    vm.editService = $.extend(true, {}, vm.service.Spec);

    vm.editMode = true;
    $timeout(function() {
      $scope.$apply();
    });
  }

  function cancelEdit() {
    vm.editMode = false;
  }

  function completeEdit() {
    vm.editService.TaskTemplate.ContainerSpec.StopGracePeriod = vm.editStopGracePeriod * 1000000000;
    if (vm.editUpdateConfigDelay) {
      vm.editService.UpdateConfig.Delay = vm.editUpdateConfigDelay * 1000000000;
    }
    vm.editService.TaskTemplate.ContainerSpec.Env = _.map(vm.editEnvVars, function(v) {
      return v.join('=');
    });

    vm.editService.Labels = _.chain(vm.editLabels)
        .keyBy('key')
        .mapValues('value')
        .value();

    vm.editService.TaskTemplate.Placement.Constraints = vm.editConstraints;

    ServiceService.update($stateParams.id, vm.editService, vm.service.Version.Index)
      .then(function(success) {
        MessageService.addSuccessMessage('Successfully updated service');
        vm.editMode = false;
        load();
      }, function(error) {
        MessageService.addErrorMessage('Error updating service', error.data);
      });
  }

  function showRemoveServiceModal() {
    $('#remove-service-modal')
      .modal({
        onApprove: function() {
          return vm.removeService($stateParams.id);
        }
      })
      .modal('show');
  }

  function showScaleServiceModal() {
    vm.numOfInstances = vm.service.Spec.Mode.Replicated.Replicas;
    $('#scale-service-modal')
      .modal({
        onApprove: function() {
          return vm.scaleService($stateParams.id, vm.numOfInstances);
        }
      })
      .modal('show');

  }

  function scaleService(id, numOfInstances) {
    ServiceService.scale(id, numOfInstances)
      .then(function(response) {
        MessageService.addSuccessMessage('Successfully scaled service');
        load();
      }, function(error) {
        MessageService.addErrorMessage('Error scaling service', error.data);
      });
  }

  function removeService(id) {
    ServiceService.remove(id)
      .then(function(response) {
        $state.go('dashboard.resources.services');
        MessageService.addSuccessMessage('Successfully removed service');
      }, function(error) {
        MessageService.addErrorMessage('Error removing service', error.data);
      });
  }

  function loadTasks() {
    var currentPage = vm.tableParams.parameters().page;
    var statusSummary = {
      inactive: 0,
      updating: 0,
      active: 0,
      errored: 0
    };

    TaskService.list()
      .then(function(tasks) {
        var instanceCount = 0;
        for (var i = 0; i < tasks.length; i++) {
          var node = vm.nodes[tasks[i].NodeID];
          if(node && node.Description) {
            tasks[i]._NodeName = vm.nodes[tasks[i].NodeID].Description.Hostname;
          }

          if (tasks[i].ServiceID === $stateParams.id) {
            instanceCount++;
            statusSummary[TaskService.stateSummaryMap[tasks[i].Status.State]]++;
          }

          if (TaskService.stateSummaryMap[tasks[i].Status.State] === 'active' || TaskService.stateSummaryMap[tasks[i].Status.State] === 'updating') {
            tasks[i]._Status = 'RUNNING';
          } else {
            tasks[i]._Status = 'HISTORICAL';
          }
        }

        vm.service._StatusSummary = statusSummary;
        vm.instanceCount = instanceCount;

        vm.tableParams.settings({
          dataset: tasks.filter(function(t) {
            return t.ServiceID === $stateParams.id;
          })
        });
        vm.tableParams.page(currentPage);

      }, function(error) {
        MessageService.addErrorMessage('Error retrieving tasks', error.data);
      });
  }

  function load() {
    ServiceService.inspect($stateParams.id)
      .then(function(service) {
          vm.service = service;
          loadTasks();
        },
        function(error) {
          MessageService.addErrorMessage('Error retrieving services', error.data);
        });
  }

}

module.exports = InspectServiceController;
