'use strict';

var _ = require('lodash');
var $ = require('jquery');
var moment = require('moment');

SwarmController.$inject = ['SwarmService', 'MessageService', '$state', '$scope'];
function SwarmController(SwarmService, MessageService, $state, $scope) {
  $scope.moment = moment;
  var vm = this;

  vm.updateConfig = updateConfig;

  vm.workerAutoAccept = false;
  vm.managerAutoAccept = false;
  vm.workerSecret = '';
  vm.managerSecret = '';


  SwarmService.getSwarmConfig()
    .then(function(data) {
      vm.swarm = data.data;
      _.forEach(vm.swarm.Spec.AcceptancePolicy.Policies, function(p) {
        if(p.Role === 'manager') {
          vm.managerAutoAccept = p.Autoaccept;
          vm.managerSecret = p.Secret;
        } else if(p.Role === 'worker') {
          vm.workerAutoAccept = p.Autoaccept;
          vm.workerSecret = p.Secret;
        }
      });
    },
    function(error) {
      MessageService.addErrorMessage('Error swarm config', error.data);
    });

  function updateConfig() {
    var updatedSwarmSpec = {};
    $.extend(updatedSwarmSpec, vm.swarm.Spec);
    updatedSwarmSpec.AcceptancePolicy.Policies = [
      {Role: 'worker', Secret: vm.workerSecret, Autoaccept: vm.workerAutoAccept},
      {Role: 'manager', Secret: vm.managerSecret, Autoaccept: vm.managerAutoAccept}
    ];
    var version = vm.swarm.Version.Index;
    SwarmService.updateSwarm(updatedSwarmSpec, version)
      .then(function(data) {
        refresh();
        MessageService.addSuccessMessage('Successfully updated swarm');
      }, function(error) {
        MessageService.addErrorMessage('Error updating swarm config', error.data.message);
      });
  }

  function refresh() {
    $state.go('dashboard.settings.swarm', {}, { reload: true });
  }
}

module.exports = SwarmController;
