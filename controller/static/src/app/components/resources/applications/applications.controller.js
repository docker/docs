'use strict';

var $ = require('jquery');
var angular = require('angular');

ApplicationsController.$inject = ['ApplicationsService', 'ContainerService', 'MessageService', '$scope', '$timeout', '$state', 'NgTableParams', 'store', '$http', '$q'];
function ApplicationsController(ApplicationsService, ContainerService, MessageService, $scope, $timeout, $state, NgTableParams, store, $http, $q) {
  /*global ORG*/
  var org = ORG || 'docker';
  /*global TAG*/
  var tag = TAG || 'latest';

  var vm = this;

  vm.tableParams = new NgTableParams({
    count: 25,
    group: '_Application'
  }, {
    groupOptions: {
      isExpanded: false
    }
  });

  vm.applications = {};
  vm.filter = '';

  vm.selected = {};
  vm.selectedApps = {};
  vm.selectedAll = false;

  vm.refresh = refresh;

  vm.showRemoveContainerDialog = showRemoveContainerDialog;
  vm.showRestartContainerDialog = showRestartContainerDialog;
  vm.showStopContainerDialog = showStopContainerDialog;
  vm.inspectContainer = inspectContainer;
  vm.removeContainer = removeContainer;
  vm.stopContainer = stopContainer;
  vm.restartContainer = restartContainer;
  vm.removeAll = removeAll;
  vm.stopAll = stopAll;
  vm.restartAll = restartAll;
  vm.checkAllApp = checkAllApp;

  vm.refresh();

  $scope.$watch(function() {
    var count = 0;
    angular.forEach(vm.selected, function (s) {
      if(s.Selected) {
        count += 1;
      }
    });
    vm.selectedItemCount = count;
  });

  // Apply jQuery to dropdowns in table once ngRepeat has finished rendering
  $scope.$on('ngRepeatFinished', function() {
    $('.ui.dropdown').dropdown();
  });

  function checkAllApp(appName) {
    angular.forEach(vm.tableParams.data, function(application) {
      angular.forEach(application.data, function(container) {
        if(container._Application === appName) {
          vm.selected[container.Id].Selected = vm.selectedApps[appName].Selected;
        }
      });
    });
  }

  function inspectContainer(containerId) {
    $state.go('dashboard.inspect.details', { id: containerId });
  }

  function restartAll() {
    var promises = [];
    for(var key in vm.selected) {
      if(vm.selected.hasOwnProperty(key) && vm.selected[key].Selected === true) {
        var promise = $http.post('/containers/' + vm.selected[key].Id + '/restart');
        promises.push(promise);
      }
    }
    $q.all(promises)
      .then(function(data) {
        deselectAll();
        vm.refresh();
        MessageService.addSuccessMessage('Restarted ' + promises.length + ' containers');
      }, function(error) {
        MessageService.addErrorMessage('Error restarting container', error.data);
      });
  }

  function stopAll() {
    var promises = [];
    for(var key in vm.selected) {
      if(vm.selected.hasOwnProperty(key) && vm.selected[key].Selected === true) {
        var promise = $http.post('/containers/' + vm.selected[key].Id + '/stop');
        promises.push(promise);
      }
    }
    $q.all(promises)
      .then(function(data) {
        deselectAll();
        vm.refresh();
        MessageService.addSuccessMessage('Stopped ' + promises.length + ' containers');
      }, function(error) {
        MessageService.addErrorMessage('Error stopping container', error.data);
      });
  }

  function removeAll() {
    var promises = [];
    for(var key in vm.selected) {
      if(vm.selected.hasOwnProperty(key) && vm.selected[key].Selected === true) {
        var promise = $http.delete('/containers/' + vm.selected[key].Id + '?v=1&force=1');
        promises.push(promise);
      }
    }
    $q.all(promises)
      .then(function(data) {
        deselectAll();
        vm.refresh();
        MessageService.addSuccessMessage('Removed ' + promises.length + ' containers');
      }, function(error) {
        MessageService.addErrorMessage('Error removing container', error.data);
      });
  }

  function deselectAll() {
    for(var id in vm.selected) {
      if(vm.selected.hasOwnProperty(id)) {
        vm.selected[id].Selected = false;
      }
    }
  }

  function showRemoveContainerDialog(container) {
    vm.selectedContainerId = container.Id;
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

  function hideRestartContainerDialog() {
    $('#restart-modal').modal('hide');
  }

  function hideStopContainerDialog() {
    $('#stop-modal').modal('hide');
  }

  function removeContainer() {
    ContainerService.remove(vm.selectedContainerId)
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
        MessageService.addErrorMessage('Error restarting container', error.data);
      });
  }

  function refresh() {
    ApplicationsService.list()
      .then(function(data) {
        var containers = [];
        vm.applications = {};

        angular.forEach(data, function(application) {
          application._ContainerCount = 0;
          application._Status = {
            'Running': 0,
            'Stopped': 0,
            'Unknown': 0
          };
          vm.selectedApps[application.name] = {Selected: vm.selectedAll};
          angular.forEach(application.services, function(service) {
            angular.forEach(service.containers, function(container) {
              vm.selected[container.Id] = {Id: container.Id, Selected: vm.selectedApps[application.name].Selected};

              var splitName = container.Names[0].split('/');
              // Container supplemental info
              container._Application = application.name;
              container._Node = nodeName(container);
              container._Name = containerName(container);
              container._Status = simpleStatus(container);
              containers.push(container);

              // Application supplemental info
              application._ContainerCount++;
              application._Status[container._Status]++;
            });
          });

          vm.applications[application.name] = application;
          vm.tableParams.settings({
            dataset: containers
          });
        });
      }, function(error) {
        MessageService.addErrorMessage('Error retrieving applications', error.data);
      });
  }

  // TODO: Move all these helper functions into the ContainersService
  function nodeName(container) {
    // Return only the node name (first component of the shortest container name)
    var components = shortestContainerName(container).split('/');
    return components[1];
  }

  function containerName(container) {
    // Remove the node name by returning the last name component of the shortest container name
    var components = shortestContainerName(container).split('/');
    return components[components.length - 1];
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

}

module.exports = ApplicationsController;
