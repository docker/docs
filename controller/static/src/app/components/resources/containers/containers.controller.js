'use strict';

var $ = require('jquery');
var angular = require('angular');

ContainersController.$inject = ['$scope', 'AuthService', 'ContainerService', 'MessageService', '$state', 'NgTableParams'];
function ContainersController($scope, AuthService, ContainerService, MessageService, $state, NgTableParams) {
  var vm = this;
  vm.account = AuthService.getCurrentUser();
  vm.containers = [];
  vm.selected = {};
  vm.selectedItemCount = 0;
  vm.selectedAll = false;
  vm.numOfInstances = 1;
  vm.selectedContainer = null;
  vm.selectedContainerId = '';
  vm.newName = '';
  vm.filter = '';

  vm.tableParams = new NgTableParams({
    count: 25,
    filter: {
      $: vm.filter,
      '_SystemContainer': false
    }
  });


  vm.runningContainersVisible = true;
  vm.toggleRunningContainers = toggleRunningContainers;
  vm.exitedContainersVisible = true;
  vm.toggleExitedContainers = toggleExitedContainers;
  vm.systemContainersVisible = false;
  vm.toggleSystemContainers = toggleSystemContainers;

  vm.showRemoveContainerDialog = showRemoveContainerDialog;
  vm.showRemoveContainersDialog = showRemoveContainersDialog;
  vm.showRestartContainerDialog = showRestartContainerDialog;
  vm.showStopContainerDialog = showStopContainerDialog;

  vm.inspectContainer = inspectContainer;
  vm.removeContainer = removeContainer;
  vm.removeVolumes = true;
  vm.stopContainer = stopContainer;
  vm.restartContainer = restartContainer;

  vm.refresh = refresh;
  vm.removeAll = removeAll;
  vm.stopAll = stopAll;
  vm.restartAll = restartAll;

  vm.refresh();

  vm.selectAll = function() {
    for(var i = 0; i < vm.tableParams.data.length; i++) {
      vm.selected[vm.tableParams.data[i].Id].Selected = vm.selectedAll;
    }
  };

  // Apply jQuery to dropdowns in table once ngRepeat has finished rendering
  $scope.$on('ngRepeatFinished', function() {
    $('.ui.container-actions.dropdown').dropdown();
    $('.circle.icon.link').popup();
  });

  // Global filtering
  $scope.$watch('vm.filter', function() {
    var filters = vm.tableParams.filter();
    filters.$ = vm.filter;
    vm.tableParams.filter(filters);
  });

  $scope.$watch('filter.$', function () {
      vm.tableParams.reload();
      vm.tableParams.page(1);
  });

  function toggleRunningContainers() {
    var filters = vm.tableParams.filter();
    if(filters._Status === 'Running') {
      return;
    } else if(filters._Status === '!Running'){
      delete filters._Status;
      vm.exitedContainersVisible = true;
      vm.runningContainersVisible = true;
    } else {
      filters._Status = '!Running';
      vm.exitedContainersVisible = true;
      vm.runningContainersVisible = false;
    }
    vm.tableParams.filter(filters);
    vm.tableParams.reload();
    vm.tableParams.page(1);
  }

  function toggleExitedContainers() {
    var filters = vm.tableParams.filter();
    if(filters._Status === 'Running') {
      delete filters._Status;
      vm.exitedContainersVisible = true;
      vm.runningContainersVisible = true;
    } else if(filters._Status === '!Running'){
      return;
    } else {
      filters._Status = 'Running';
      vm.exitedContainersVisible = false;
      vm.runningContainersVisible = true;
    }
    vm.tableParams.filter(filters);
    vm.tableParams.reload();
    vm.tableParams.page(1);
  }

  function toggleSystemContainers() {
    var filters = vm.tableParams.filter();
    if(filters._SystemContainer === false) {
      delete filters._SystemContainer;
      vm.systemContainersVisible = true;
    } else {
      filters._SystemContainer = false;
      vm.systemContainersVisible = false;
    }

    vm.tableParams.filter(filters);
    vm.tableParams.reload();
    vm.tableParams.page(1);
  }

  function updateSelectedAllCheckbox() {
    if(vm.tableParams.data.length === 0) {
      vm.selectedAll = false;
      return;
    }

    var selectedAll = true;
    for(var i = 0; i < vm.tableParams.data.length; i++) {
      if(!vm.selected[vm.tableParams.data[i].Id].Selected) {
        selectedAll = false;
        break;
      }
    }
    vm.selectedAll = selectedAll;
  }

  // If table data changes through page number or page size changes, update selections
  $scope.$watch('vm.tableParams.data', function() {
    // Ensure only visible items remain selected
    Object.keys(vm.selected).forEach(function(key, index) {
      if(vm.selected[key].Selected === true) {
        for(var i = 0; i < vm.tableParams.data.length; i++) {
          if(vm.tableParams.data[i].Id === key) {
            return;
          }
        }
        vm.selected[key].Selected = false;
      }
    });

    // If the page number or page size changes, we need to update the selected all checkbox
    updateSelectedAllCheckbox();
  });

  // If items are selected, refresh counts and status of selected all checkbox
  $scope.$watch('vm.selected', function() {
    // Update selected count
    var count = 0;
    angular.forEach(vm.selected, function (s) {
      if(s.Selected) {
        count += 1;
      }
    });
    vm.selectedItemCount = count;

    updateSelectedAllCheckbox();
  }, true);

  function tableFilter(row) {
    return (
      angular.lowercase(row._Status).indexOf(angular.lowercase(vm.filterQuery) || '') !== -1 ||
        angular.lowercase(row.Id).indexOf(angular.lowercase(vm.filterQuery) || '') !== -1 ||
        angular.lowercase(row._Name).indexOf(angular.lowercase(vm.filterQuery) || '') !== -1 ||
        angular.lowercase(row.Image).indexOf(angular.lowercase(vm.filterQuery) || '') !== -1 ||
        angular.lowercase(row._Node).indexOf(angular.lowercase(vm.filterQuery) || '') !== -1
    );
  }

  function inspectContainer(containerId) {
    $state.go('dashboard.inspect.details', { id: containerId });
  }

  function restartAll() {
    // TODO: Add success message for multi-container operations, we probably need to construct some
    // sort of promise that only shows a success message after all operations have finished.

    MessageService.addSuccessMessage('Restarting selected containers');
    angular.forEach(vm.selected, function (s) {
      if(s.Selected === true) {
        ContainerService.restart(s.Id)
        .then(function(data) {
          delete vm.selected[s.Id];
          vm.refresh();
        }, function(error) {
          MessageService.addErrorMessage('Error restarting container', error.data);
        });
      }
    });
  }

  function stopAll() {
    // TODO: Add success message for multi-container operations, we probably need to construct some
    // sort of promise that only shows a success message after all operations have finished.

    MessageService.addSuccessMessage('Stopping selected containers');
    angular.forEach(vm.selected, function (s) {
      if(s.Selected === true) {
        ContainerService.stop(s.Id)
          .then(function(data) {
            delete vm.selected[s.Id];
            vm.refresh();
          }, function(error) {
            MessageService.addErrorMessage('Error stopping container', error.data);
          });
      }
    });
  }

  function removeAll() {
    // TODO: Add success message for multi-container operations, we probably need to construct some
    // sort of promise that only shows a success message after all operations have finished.

    MessageService.addSuccessMessage('Removing selected containers');
    angular.forEach(vm.selected, function (s) {
      if(s.Selected === true) {
        ContainerService.remove(s.Id)
          .then(function(data) {
            delete vm.selected[s.Id];
            vm.refresh();
          }, function(error) {
            MessageService.addErrorMessage('Error removing container', error.data);
          });
      }
    });
    hideRemoveContainersDialog();
  }

  function nodeName(container) {
    // Return only the node name (first component of the shortest container name)
    var components = shortestContainerName(container).split('/');
    return components[1];
  }

  function shortestContainerName(container) {
    // Distill shortest container name by taking the name with the fewest components
    // Names with the same number of components are considered in undefined order
    var shortestName = '';
    var minComponents = 99;

    var names = container.Names;
    for(var i = 0; i < names.length; i++) {
      var name = names[i];
      var numComponents = name.split('/').length;
      if(numComponents < minComponents) {
        shortestName = name;
        minComponents = numComponents;
      }
    }

    return shortestName;
  }

  function containerName(container) {
    // Remove the node name by returning the last name component of the shortest container name
    var components = shortestContainerName(container).split('/');
    return components[components.length - 1];
  }

  function showRemoveContainerDialog(container) {
    vm.selectedContainerId = container.Id;
    $('#remove-modal').modal('show');
  }

  function showRemoveContainersDialog() {
    vm.selectedContainerId = null;
    $('#remove-modal').modal('show');
  }

  function showRestartContainerDialog(container) {
    vm.selectedContainerId = container.Id;
    $('#restart-modal').modal('show');
  }

  function showStopContainerDialog(container) {
    vm.selectedContainerId = container.Id;
    $('#stop-modal').modal('show');
  }

  function hideRemoveContainerDialog() {
    $('#remove-modal').modal('hide');
  }

  function hideRemoveContainersDialog() {
    $('#remove-modal').modal('hide');
  }

  function hideRestartContainerDialog() {
    $('#restart-modal').modal('hide');
  }

  function hideStopContainerDialog() {
    $('#stop-modal').modal('hide');
  }

  function removeContainer() {
    ContainerService.remove(vm.selectedContainerId, vm.removeVolumes)
      .then(function(data) {
        MessageService.addSuccessMessage('Success', 'Removed container ' + vm.selectedContainerId);
        vm.refresh();
        hideRemoveContainerDialog();
      }, function(error) {
        MessageService.addErrorMessage('Error removing container', error.data);
      });
  }

  function stopContainer() {
    ContainerService.stop(vm.selectedContainerId)
      .then(function(data) {
        MessageService.addSuccessMessage('Success', 'Stopped container ' + vm.selectedContainerId);
        vm.refresh();
        hideStopContainerDialog();
      }, function(error) {
        MessageService.addErrorMessage('Error stopping container', error.data);
      });
  }

  function restartContainer() {
    ContainerService.restart(vm.selectedContainerId)
      .then(function(data) {
        MessageService.addSuccessMessage('Success', 'Restarted container ' + vm.selectedContainerId);
        vm.refresh();
        hideRestartContainerDialog();
      }, function(error) {
        MessageService.addErrorMessage('Error stopping container', error.data);
      });
  }

  function simpleStatus(container) {
    if(container.Status.indexOf('Paused') !== -1){
      return 'Paused';
    }
    if(container.Status.indexOf('Up') === 0){
      return 'Running';
    }
    else if(container.Status.indexOf('Exited') === 0){
      return 'Stopped';
    }
    return 'Unknown';
  }

  function refresh() {
    ContainerService.list(true, false)
      .then(function(data) {
        vm.containers = data;
        angular.forEach(vm.containers, function (container) {
          vm.selected[container.Id] = {Id: container.Id, Selected: vm.selectedAll};
          container._Node = nodeName(container);
          container._Name = containerName(container);
          container._Status = simpleStatus(container);
          container._Label = container.Labels['com.docker.ucp.access.label'];
          container._SystemContainer = container.Labels['com.docker.ucp.InstanceID'] ? true : false;
        });

        vm.tableParams.settings({
          dataset: vm.containers
        });
      }, function(error) {
        MessageService.addErrorMessage('Error retrieving container information', error.data);
      });

    vm.containers = [];
    vm.selected = {};
    vm.selectedItemCount = 0;
    vm.selectedAll = false;
    vm.numOfInstances = 1;
    vm.selectedContainerId = '';
    vm.newName = '';
  }

}

module.exports = ContainersController;
