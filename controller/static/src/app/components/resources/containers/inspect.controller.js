'use strict';

var $ = require('jquery');

InspectDetailsController.$inject = ['resolvedContainer', 'ContainerService', 'MessageService', 'NgTableParams'];
function InspectDetailsController(resolvedContainer, ContainerService, MessageService, NgTableParams) {
  var vm = this;
  vm.healthLog = new NgTableParams({}, {});
  vm.container = resolvedContainer;
  vm.showContainerDialog = showContainerDialog;
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

  function updateHealthLog() {
    if(vm.container.State.Health) {
      vm.healthLog.settings({
        dataset: vm.container.State.Health.Log
      });
    }
  }

  function refresh() {
    ContainerService
      .inspect(vm.container.Id)
      .then(function(data) {
        vm.container = data;
        updateHealthLog();
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

  function showContainerDialog(className) {
    $('.ui.small.' + className + '.modal').modal('show');
  }

  updateHealthLog();

  // Semantic UI functions
  $('.menu .item').tab();
}

module.exports = InspectDetailsController;
