'use strict';

var $ = require('jquery');
var _ = require('lodash');

NetworksController.$inject = ['NetworksService', 'MessageService', '$state', '$timeout', '$scope', 'NgTableParams'];
function NetworksController(NetworksService, MessageService, $state, $timeout, $scope, NgTableParams) {
  var vm = this;
  vm.tableParams = new NgTableParams({
    count: 25
  }, {});
  vm.selectedNetwork = null;
  vm.networkName = '';
  vm.networkDriver = '';
  vm.networkOptions = '';
  vm.networkIPAMDriver = '';
  vm.networkIPAMSubnet = '';
  vm.networkIPAMIPRange = '';
  vm.networkIPAMGateway = '';
  vm.filter = '';
  vm.networkLabels = [];
  vm.hasLabels = false;
  vm.hasOptions = false;
  vm.labelName = '';
  vm.labelValue = '';
  vm.accessLabel = '';

  vm.refresh = refresh;
  vm.createNetwork = createNetwork;
  vm.removeNetwork = removeNetwork;
  vm.showCreateNetworkDialog = showCreateNetworkDialog;
  vm.showRemoveNetworkDialog = showRemoveNetworkDialog;
  vm.showNetworkDialog = showNetworkDialog;
  vm.pushLabel = pushLabel;

  refresh();

  function pushLabel() {
    var label = {
      name: vm.labelName,
      value: vm.labelValue
    };
    vm.networkLabels.push(label);
    vm.labelName = '';
    vm.labelValue = '';
  }

  function showCreateNetworkDialog() {
    $('#create-network-modal')
      .modal({
        onApprove: function() {
          return vm.createNetwork();
        }
      })
      .modal('refresh')
      .modal('show');
  }

  function hideCreateNetworkDialog() {
    $('#create-network-modal')
      .modal('hide');
  }

  function showNetworkDialog(network) {
    vm.selectedNetwork = network;
    vm.hasLabels = !_.isEmpty(vm.selectedNetwork.Labels);
    vm.hasOptions = !_.isEmpty(vm.selectedNetwork.Options);
    $('#show-network-modal').modal('show');
  }

  function showRemoveNetworkDialog(network) {
    vm.selectedNetwork = network;
    $('#remove-network-modal').modal('show');
  }

  // Global filtering
  $scope.$watch('vm.filter', function() {
    vm.tableParams.filter({ $: vm.filter });
  });

  $scope.$watch('filter.$', function () {
    vm.tableParams.reload();
    vm.tableParams.page(1);
  });

  function refresh() {
    // reset create net fields
    vm.networkName = '';
    vm.networkDriver = '';
    vm.networkOptions = '';
    vm.networkIPAMDriver = 'default';
    vm.networkIPAMSubnet = '';
    vm.networkIPAMIPRange = '';
    vm.networkIPAMGateway = '';
    vm.accessLabel = '';

    NetworksService.list()
      .then(function(data) {
        vm.tableParams.settings({dataset: data});
      }, function(error) {
        MessageService.addErrorMessage('Error Retrieving Networks', error.data);
      });
  }

  function createNetwork() {
    vm.createNetworkError = '';

    if(!$('#create-network-modal .ui.form').form('validate form')) {
      return false;
    }

    var config;
    if (vm.networkIPAMSubnet === '' && vm.networkIPAMIPRange === '' && vm.networkIPAMGateway === '') {
      config = null;
    } else {
      config = [{
        Subnet: vm.networkIPAMSubnet !== '' ? vm.networkIPAMSubnet : null,
        IPRange: vm.networkIPAMIPRange !== '' ? vm.networkIPAMIPRange : null,
        Gateway: vm.networkIPAMGateway !== '' ? vm.networkIPAMGateway : null
      }];
    }

    var networkOptions = {};
    if (vm.networkOptions !== '') {
      var opts = vm.networkOptions.split(' ');
      for (var i = 0; i < opts.length; i++) {
        var opt = opts[i].split('=');
        networkOptions[opt[0]] = opt[1];
      }
    }

    var labels = {};
    if (vm.labelName && vm.labelName.length > 0) {
      labels[vm.labelName] = vm.labelValue;
    }

    _.forEach(vm.networkLabels, function(label, index) {
      labels[label.name] = label.value;
    });

    if (vm.accessLabel !== '') {
      labels['com.docker.ucp.access.label'] = vm.accessLabel;
    }

    var payload = {
      Name: vm.networkName,
      Driver: vm.networkDriver,
      CheckDuplicate: true,
      Options: networkOptions,
      IPAM: {
        Driver: vm.networkIPAMDriver,
        Config: config
      },
      Labels: labels
    };

    NetworksService.create(payload)
      .then(function(data) {
        MessageService.addSuccessMessage('Success', 'Created network');
        vm.refresh();
        hideCreateNetworkDialog();
      }, function(error) {
        vm.createNetworkError = error.data;
        showCreateNetworkDialog();
      });

    return true;
  }

  function removeNetwork() {
    NetworksService.remove(vm.selectedNetwork)
      .then(function(data) {
        MessageService.addSuccessMessage('Success', 'Removed network');
        vm.refresh();
      }, function(error) {
        MessageService.addErrorMessage('Error removing network', error.data);
      });
  }
}

module.exports = NetworksController;
