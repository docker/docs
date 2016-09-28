'use strict';

var $ = require('jquery');
var _ = require('lodash');

TeamPermissionsController.$inject = ['team', 'permissions', '$scope', '$state', 'AccountsService', 'MessageService', 'NgTableParams'];
function TeamPermissionsController(team, permissions, $scope, $state, AccountsService, MessageService, NgTableParams) {
  var vm = this;
  vm.team = team;
  vm.permissions = permissions;
  vm.tableParams = new NgTableParams({}, {
    dataset: permissions
  });

  vm.selectedPermission = {};
  vm.addLabelVisible = false;
  vm.newLabelRequest = {
    teamId: vm.team.id,
    role: 0
  };

  vm.showAddLabel = showAddLabel;
  vm.resetFormVisibility = resetFormVisibility;
  vm.addLabel = addLabel;
  vm.deletePermission = deletePermission;
  vm.showDeletePermissionModal = showDeletePermissionModal;
  vm.roleDescription = roleDescription;

  function roleDescription(roleId) {
    if(roleId === 0) {
      return 'No Access';
    } else if(roleId === 1) {
      return 'View Only';
    } else if(roleId === 2) {
      return 'Restricted Control';
    } else if(roleId === 3) {
      return 'Full Control';
    }
    return 'Unknown';
  }

  function showAddLabel() {
    //vm.addLabelVisible = true;
    $('#add-label-modal').modal('show');
  }

  function showDeletePermissionModal(permission) {
    vm.selectedPermission = permission;
    $('#delete-permission-modal').modal('show');
  }

  function resetFormVisibility() {
    vm.addLabelVisible = false;
  }

  $scope.$watch('vm.filter', function() {
    vm.tableParams.filter({ $: vm.filter });
  });

  $scope.$watch('filter.$', function () {
      vm.tableParams.reload();
      vm.tableParams.page(1);
  });

  function refresh() {
    $state.go('dashboard.accounts.teamPermissions', { teamId: vm.team.id }, { reload: true });
  }

  function deletePermission() {
    AccountsService.deletePermission(vm.selectedPermission)
      .then(function(data) {
        refresh();

        MessageService.addSuccessMessage('Success', 'Removed permission label ' + vm.selectedPermission.label);
      }, function(error) {
        MessageService.addErrorMessage('Error removing permission label', error.data);
      });
  }

  function addLabel() {

    if(_.find(vm.permissions, {label: vm.newLabelRequest.label})) {
      MessageService.addErrorMessage('Label already exists', 'A label with the name ' + vm.newLabelRequest.label + ' already exists');
      return;
    }

    AccountsService.addPermission(vm.newLabelRequest)
      .then(function(data) {
        refresh();

        MessageService.addSuccessMessage('Success', 'Added permission label ' + vm.newLabelRequest.label);
      }, function(error) {
        MessageService.addErrorMessage('Error adding permission label', error.data);
      });
  }
}

module.exports = TeamPermissionsController;
