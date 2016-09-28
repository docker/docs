'use strict';

var $ = require('jquery');

InspectBaseController.$inject = ['resolvedContainer', 'ContainerService', 'MessageService', '$state', '$rootScope'];
function InspectBaseController(resolvedContainer, ContainerService, MessageService, $state, $rootScope) {
  if($rootScope.fromState && $rootScope.fromState.name !== '') {
    $state.$current.parent.self.ncyBreadcrumb.parent = $rootScope.fromState.name + '(' + JSON.stringify($rootScope.fromParams) + ')';
  }
  $state.$current.parent.self.ncyBreadcrumb.label = 'Container: ' + resolvedContainer.Name.split('/').pop();

  var vm = this;
  vm.container = resolvedContainer;
  vm.showContainerDialog = showContainerDialog;
  vm.restartContainer = restartContainer;
  vm.stopContainer = stopContainer;
  vm.removeContainer = removeContainer;
  vm.removeVolumes = true;
  vm.renameContainer = renameContainer;
  vm.parseLinkingString = parseLinkingString;
  vm.isEmptyObject = isEmptyObject;
  vm.links = parseContainerLinks(vm.container.HostConfig.Links);
  vm.labels = vm.container.Config.Labels;
  vm.changeTabIfRunning = changeTabIfRunning;

  if(vm.container.State.Running) {
    ContainerService.top(vm.container.Id)
      .then(function(data) {
        vm.top = data;
      }, function(error) {
        MessageService.addErrorMessage('Error retrieving process information', error.data);
      });
  }

  function refresh() {
    ContainerService
      .inspect(vm.container.Id)
      .then(function(data) {
        vm.container = data;
      }, function(error) {
        MessageService.addErrorMessage('Error retrieving container details', error.data);
      });
  }

  function changeTabIfRunning(tabName) {
    if(!vm.container.State.Running) {
      return;
    }

    $.tab('change tab', tabName);
  }

  function parseContainerLinks(links) {
    var l = [];
    if (links === null) {
      return l;
    }
    for (var i = 0; i < links.length; i++) {
      var parts = links[i].split(':');
      var link = {
        container: parts[0].slice(1),
        link: parts[1].split('/')[2]
      };
      l.push(link);
    }

    return l;
  }

  function parseLinkingString(linkingString) {
    var linkedTo = linkingString.split(':')[0].replace('/', '');
    var alias = linkingString.split(':')[1];

    return linkedTo + String.fromCharCode(8594) + alias.substring(alias.lastIndexOf('/') + 1, alias.length);
  }

  function isEmptyObject(obj) {
    for (var k in obj) {
      return false;
    }

    return true;
  }

  function showContainerDialog(id) {
    $(id).modal('show');
  }

  function restartContainer() {
    ContainerService.restart(vm.container.Id)
      .then(function(data) {
        MessageService.addSuccessMessage('Success', 'Restarted container ' + vm.container.Id);
        refresh();
      }, function(error) {
        MessageService.addErrorMessage('Error restarting container', error.data);
      });
  }

  function stopContainer() {
    ContainerService.stop(vm.container.Id)
      .then(function(data) {
        MessageService.addSuccessMessage('Success', 'Stopped container ' + vm.container.Id);
        refresh();
      }, function(error) {
        MessageService.addErrorMessage('Error stopping container', error.data);
      });
  }

  function removeContainer() {
    ContainerService.remove(vm.container.Id, vm.removeVolumes)
      .then(function(data) {
        $state.go('dashboard.resources.containers');
        MessageService.addSuccessMessage('Success', 'Removed container ' + vm.container.Id);
      }, function(error) {
        MessageService.addErrorMessage('Error removing container', error.data);
      });
  }

  function renameContainer() {
    ContainerService.rename(vm.container.Id, vm.newName)
      .then(function(data) {
        MessageService.addSuccessMessage('Success', 'Renamed container ' + vm.container.Id + ' to ' + vm.newName);
        refresh();
      }, function(error) {
        MessageService.addErrorMessage('Error renaming container', error.data);
      });
  }

  // Semantic UI functions
  $('.menu .item').tab();
}

module.exports = InspectBaseController;
