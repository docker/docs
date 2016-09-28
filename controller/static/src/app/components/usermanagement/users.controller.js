'use strict';

var $ = require('jquery');

UsersController.$inject = ['users', '$scope', '$timeout', '$http', '$state', 'AccountsService', 'MessageService', 'NgTableParams'];
function UsersController(users, $scope, $timeout, $http, $state, AccountsService, MessageService, NgTableParams) {
  var vm = this;
  vm.filter = '';
  vm.tableParams = new NgTableParams({}, {
    dataset: users
  });

  vm.createUserVisible = false;
  vm.newUserRequest = {};
  vm.editedAccount = {};
  vm.adminToggled = true;

  vm.createUser = createUser;
  vm.showCreateUser = showCreateUser;
  vm.removeAccount = removeAccount;
  vm.showRemoveAccountDialog = showRemoveAccountDialog;

  vm.resetFormVisibility = resetFormVisibility;
  vm.refresh = refresh;

  vm.showEditUser = showEditUser;
  vm.editAccount = editAccount;

  function isValid() {
    return $('#editAccount .ui.form').form('validate form');
  }

  function showEditUser(account) {
    vm.adminToggled = false;
    vm.editedAccount = account;
    $('#edit-user-modal')
      .modal({
        onApprove: function() {
          return vm.editAccount();
        }
      })
      .modal('show');
  }

  function editAccount() {
    // If a new password has been entered, then update it
    if(vm.newPassword !== '') {
      vm.editedAccount.password = vm.newPassword;
    }

    $http.post('/api/accounts', vm.editedAccount)
      .success(function(data) {
        MessageService.addSuccessMessage('Success', 'Updated ' + vm.editedAccount.username);
        $state.go('dashboard.accounts.users');
      })
      .error(function(error) {
        MessageService.addErrorMessage('Error updating user', error.data);
      });
  }

  function refresh() {
    AccountsService.list()
      .then(function(data) {
        vm.tableParams.settings({dataset: data});
      }, function(error) {
        MessageService.addErrorMessage('Error retrieving users', error);
      });
  }

  function showRemoveAccountDialog(account) {
    vm.selectedAccount = account;
    $('#remove-modal').modal('show');
  }

  function hideRemoveAccountDialog() {
    $('#remove-modal').modal('hide');
  }

  function removeAccount() {
    AccountsService.removeAccount(vm.selectedAccount)
      .then(function(data) {
        MessageService.addSuccessMessage('Success', 'Removed user ' + vm.selectedAccount.username);
        vm.refresh();
        hideRemoveAccountDialog();
      }, function(error) {
        MessageService.addErrorMessage('Error removing user', error.data);
      });
  }

  function resetFormVisibility() {
    vm.newUserRequest = {
      role: 0,
      admin: false
    };
    vm.createUserVisible = false;
    $timeout(function() {
      $('.ui.default-permissions.dropdown').dropdown('set selected', 'No Access');
    });
  }

  function showCreateUser() {
    resetFormVisibility();
    $('#create-user-modal')
      .modal({
        onApprove: function() {
          return vm.createUser();
        }
      })
      .modal('show');
  }

  function hideCreateUserDialog() {
    $('#create-user-modal').modal('hide');
  }

  function createUser() {
    if (!$('.ui.new-user.form').form('validate form')) {
      return false;
    }

    AccountsService.createUser(vm.newUserRequest)
      .then(function(data) {
        MessageService.addSuccessMessage('Success', 'Created user ' + vm.newUserRequest.username);
        resetFormVisibility();
        vm.refresh();
        hideCreateUserDialog();
      }, function(error) {
        MessageService.addErrorMessage('Error creating user', error.data);
      });
    return true;
  }

  $scope.$watch('vm.editedAccount.admin', function(newValue, oldValue) {
    $timeout(function() {
      if(oldValue !== undefined && newValue !== oldValue) {
        vm.adminToggled = true;
      }

      if(newValue) {
        $('.ui.default-permissions.dropdown').addClass('disabled');
        $('.ui.default-permissions.dropdown').dropdown('set selected', 'Full Control');
      } else {
        $('.ui.default-permissions.dropdown').removeClass('disabled');
        $('.ui.default-permissions.dropdown').dropdown('set selected', 'No Access');
      }
    });
  });

  $scope.$watch('vm.newUserRequest.admin', function(isAdmin) {
    $timeout(function() {
      if(isAdmin) {
        $('.ui.default-permissions.dropdown').addClass('disabled');
        $('.ui.default-permissions.dropdown').dropdown('set selected', 'Full Control');
      } else {
        $('.ui.default-permissions.dropdown').removeClass('disabled');
        $('.ui.default-permissions.dropdown').dropdown('set selected', 'No Access');
      }
    });
  });

  $scope.$watch('vm.filter', function() {
    vm.tableParams.filter({ $: vm.filter });
  });

  $scope.$watch('filter.$', function () {
      vm.tableParams.reload();
      vm.tableParams.page(1);
  });

  // Set the team selector to 'All Users'
  $timeout(function() {
    $('.select-team-widget .ui.dropdown').dropdown('set selected', 'All Users');
    $scope.$apply();
  });
}

module.exports = UsersController;
