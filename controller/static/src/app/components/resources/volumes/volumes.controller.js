'use strict';

var $ = require('jquery');
var angular = require('angular');

VolumesController.$inject = ['VolumesService', 'MessageService', '$state', '$timeout', '$scope', 'NgTableParams'];
function VolumesController(VolumesService, MessageService, $state, $timeout, $scope, NgTableParams) {
  var vm = this;
  vm.tableParams = new NgTableParams({
    count: 25
  }, {});
  vm.selected = {};
  vm.selectedItemCount = 0;
  vm.selectedAll = false;
  vm.selectedVolume = null;
  vm.refresh = refresh;
  vm.createVolume = createVolume;
  vm.removeVolume = removeVolume;
  vm.showCreateVolumeDialog = showCreateVolumeDialog;
  vm.showRemoveVolumeDialog = showRemoveVolumeDialog;
  vm.removeAll = removeAll;

  vm.volumeName = '';
  vm.volumeDriver = '';
  vm.volumeOptions = '';

  refresh();

  vm.selectAll = function() {
    for(var i = 0; i < vm.tableParams.data.length; i++) {
      vm.selected[vm.tableParams.data[i].Name].Selected = vm.selectedAll;
    }
  };

  function updateSelectedAllCheckbox() {
    if(vm.tableParams.data.length === 0) {
      vm.selectedAll = false;
      return;
    }

    var selectedAll = true;
    for(var i = 0; i < vm.tableParams.data.length; i++) {
      if(!vm.selected[vm.tableParams.data[i].Name].Selected) {
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
          if(vm.tableParams.data[i].Name === key) {
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

  function removeAll() {
    MessageService.addSuccessMessage('Removing selected volumes');
    angular.forEach(vm.selected, function (s) {
      if(s.Selected === true) {
        VolumesService.remove(s.Name)
          .then(function(data) {
            delete vm.selected[s.Name];
            vm.refresh();
          }, function(error) {
            MessageService.addErrorMessage('Error removing volume', error.data);
          });
      }
    });
    hideRemoveVolumeDialog();
  }

  function showCreateVolumeDialog(volume) {
    $('#create-volume-modal')
      .modal({
        onApprove: function() {
          return vm.createVolume();
        }
      })
      .modal('show');
  }

  function hideRemoveVolumeDialog() {
    $('#remove-volume-modal').modal('hide');
  }

  function hideCreateVolumeDialog() {
    $('#create-volume-modal').modal('hide');
  }

  function showRemoveVolumeDialog(volume) {
    vm.selectedVolume = volume;
    $('#remove-volume-modal').modal('show');
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
    vm.selectedVolume = null;
    vm.volumeName = '';
    vm.volumeDriver = '';
    vm.volumeOptions = '';
    VolumesService.list()
      .then(function(data) {
        angular.forEach(data, function (v) {
          vm.selected[v.Name] = {Name: v.Name, Selected: vm.selectedAll};
        });
        vm.tableParams.settings({dataset: data});
      }, function(data) {
        MessageService.addErrorMessage('Error retrieving volumes', data.data);
      });
  }

  function createVolume() {
    if(!$('#create-volume-modal .ui.form').form('validate form')) {
        return false;
    }

    var volumeOptions = {};
    if (vm.volumeOptions !== '') {
      var opts = vm.volumeOptions.split(' ');
      for (var i = 0; i < opts.length; i++) {
        var opt = opts[i].split('=');
        volumeOptions[opt[0]] = opt[1];
      }
    }

    var payload = {
      Name: vm.volumeName,
      Driver: vm.volumeDriver,
      DriverOpts: volumeOptions
    };

    VolumesService.create(payload)
      .then(function(data) {
        MessageService.addSuccessMessage('Volume Created', 'Successfully created volume ' + vm.volumeName);
        vm.refresh();
        hideCreateVolumeDialog();
      }, function(data) {
        MessageService.addErrorMessage('Error creating volume', data.data);
      });

    return true;
  }

  function removeVolume() {
    VolumesService.remove(vm.selectedVolume.Name)
      .then(function(data) {
        MessageService.addSuccessMessage('Volume Removed', 'Successfully removed volume');
        vm.refresh();
        hideRemoveVolumeDialog();
      }, function(data) {
        MessageService.addErrorMessage('Error removing volume', data.data);
      });
  }
}

module.exports = VolumesController;
