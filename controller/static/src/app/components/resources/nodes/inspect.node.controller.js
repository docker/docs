'use strict';

var _ = require('lodash');
var $ = require('jquery');

InspectNodeController.$inject = ['LogService', 'NodeService', 'MessageService', 'TaskService', 'NgTableParams', '$state', '$stateParams', '$timeout', '$scope'];
function InspectNodeController(LogService, NodeService, MessageService, TaskService, NgTableParams, $state, $stateParams, $timeout, $scope) {
  var websocket;
  var id = $stateParams.id;

  var vm = this;
  vm.node = {};
  vm.tableParams = new NgTableParams({
    filter: {
      '_Status': 'RUNNING'
    },
    sorting: {
      'Status.State': 'desc',
      Instance: 'asc'
    }
  });
  vm.stateSummaryMap = TaskService.stateSummaryMap;
  vm.editNodeLabels = [];
  vm.runningTasksVisible = true;
  vm.historicalTasksVisible = false;
  vm.activeTab = 'details';
  vm.agentTask = null;
  vm.nodeLogs = '';

  vm.inspectContainer = inspectContainer;
  vm.removeNodeLabel = removeNodeLabel;
  vm.startEdit = startEdit;
  vm.cancelEdit = cancelEdit;
  vm.completeEdit = completeEdit;
  vm.toggleRunningTasks = toggleRunningTasks;
  vm.toggleHistoricalTasks = toggleHistoricalTasks;
  vm.showLogs = showLogs;
  vm.showDetails = showDetails;
  vm.showTasks = showTasks;
  vm.removeNode = removeNode;

  function showTasks() {
    vm.activeTab = 'tasks';
  }

  function showDetails() {
    vm.activeTab = 'details';
  }

  function showLogs() {
    if(!vm.agentTask) {
      // If there are no logs available, do nothing
      return;
    }

    vm.nodeLogs = '';
    vm.activeTab = 'logs';

    LogService.get(vm.agentTask.Status.ContainerStatus.ContainerID)
      .then(function(ws) {
        websocket = ws;
        websocket.onopen = function(evt) {


          websocket.onmessage = function(msg) {
            vm.nodeLogs += msg.data;
            $timeout(function() {
              $scope.$apply();
            });
          };
        };
      }, function(error) {
        MessageService.addErrorMessage('Could not retrieve logs');
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

  function inspectContainer(containerId) {
    if (!containerId) {
      return;
    }

    $state.go('dashboard.inspect.details', {
      id: containerId
    });
  }

  function startEdit() {
    vm.editNodeLabels = _.map(vm.node.Spec.Labels, function (v, k) {
      return { id: k, value: v };
    });
    load();
    vm.editMode = true;
  }

  function cancelEdit() {
    load();
    vm.editMode = false;
  }

  function completeEdit() {
    var labelMap = {};
    for(var i = 0; i < vm.editNodeLabels.length; i++) {
      labelMap[vm.editNodeLabels[i].id] = vm.editNodeLabels[i].value;
    }

    vm.node.Spec.Labels = labelMap;

    updateNode(vm.node.Spec)
      .then(function(r) {
        vm.editMode = false;
      });
  }

  function updateNode(spec) {
    return NodeService.update(vm.node.ID, vm.node.Version.Index, spec)
      .then(function(response) {
        load();
      }, function(error) {
        MessageService.addErrorMessage('Error updating node', error.data.message);
      });
  }

  function removeNodeLabel(k) {
    vm.editNodeLabels = _.filter(vm.editNodeLabels, function(l) {
      return l.id !== k;
    });
  }

  function filterTasksForAgent(t) {
    return t.NodeID === id && t.Spec.ContainerSpec.Image.indexOf('ucp-agent') > -1;
  }

  function removeNode() {
    NodeService.remove($stateParams.id)
      .then(function(success) {
        $state.go('dashboard.resources.nodes');
        MessageService.addSuccessMessage('Successfully removed node');
      }, function(error) {
        MessageService.addErrorMessage('Error removing node', error.data.message);
      });
  }


  function load() {
    NodeService.inspect($stateParams.id)
      .then(function(node) {
        vm.node = node;
        TaskService.listForNode(vm.node.ID)
          .then(function(tasks) {
            vm.agentTask = _.find(tasks, filterTasksForAgent);
            vm.tableParams.settings({
              dataset: tasks
            });
          }, function(error) {
            MessageService.addErrorMessage('Error retrieving tasks', error.data);
          });
      }, function(error) {
        MessageService.addErrorMessage('Error retrieving node', error.data);
      });
  }

  load();
}

module.exports = InspectNodeController;
