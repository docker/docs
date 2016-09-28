'use strict';

var $ = require('jquery');

TeamController.$inject = ['team', 'teamMembers', 'users', '$scope', '$state', 'AccountsService', 'MessageService', 'NgTableParams'];
function TeamController(team, teamMembers, users, $scope, $state, AccountsService, MessageService, NgTableParams) {
  $state.current.ncyBreadcrumb.label = 'Team: ' + team.name;

  var vm = this;
  vm.team = team;
  vm.teamMembers = teamMembers;
  vm.users = users;

  vm.filter = '';
  vm.addUserToTeamFilter = '';

  vm.tableParams = new NgTableParams({}, {
    dataset: vm.teamMembers
  });

  vm.addUserTableParams = new NgTableParams({}, {
    dataset: vm.users
  });

  vm.addUserVisible = false;

  vm.showRemoveFromTeamModal = showRemoveFromTeamModal;
  vm.removeFromTeam = removeFromTeam;
  vm.showAddUser = showAddUser;
  vm.addUserToTeam = addUserToTeam;
  vm.isUserInSelectedTeam = isUserInSelectedTeam;
  vm.resetFormVisibility = resetFormVisibility;
  vm.refresh = refresh;

  function resetFormVisibility() {
    vm.addUserToTeamFilter = '';
    vm.addUserVisible = false;
  }

  function showRemoveFromTeamModal(account) {
    vm.selectedAccount = account;
    $('#remove-from-team-modal').modal('show');
  }

  function removeFromTeam() {
    AccountsService.removeUserFromTeam(vm.team.id, vm.selectedAccount.username)
      .then(function(data) {
        MessageService.addSuccessMessage('Success', 'Removed user ' + vm.selectedAccount.username + ' from ' + vm.team.name);
        vm.refresh();
      }, function(error) {
        MessageService.addErrorMessage('Error removing user from team', error.data);
      });
  }

  function showAddUser() {
    resetFormVisibility();
    /*vm.addUserVisible = true;*/
    $('#add-user-to-team-modal').modal('show');
  }

  function isUserInSelectedTeam(username) {
    if(vm.team.managed_members && $.inArray(username, vm.team.managed_members) >= 0) {
      return true;
    }
    if(vm.team.discovered_members && $.inArray(username, vm.team.discovered_members) >= 0) {
      return true;
    }
    return false;
  }

  function addUserToTeam(username) {
    AccountsService.addUserToTeam(vm.team.id, username)
      .then(function(data) {
        MessageService.addSuccessMessage('Success', 'Added user ' + username + ' to ' + vm.team.name);
        vm.refresh();
      }, function(error) {
        MessageService.addErrorMessage('Error adding user to team', error.data);
      });
  }

  $scope.$watch('vm.filter', function() {
    vm.tableParams.filter({ $: vm.filter });
  });

  $scope.$watch('vm.addUserToTeamFilter', function() {
    vm.addUserTableParams.filter({ $: vm.addUserToTeamFilter });
  });

  $scope.$watch('filter.$', function () {
      vm.tableParams.reload();
      vm.tableParams.page(1);
      vm.addUserTableParams.reload();
      vm.addUserTableParams.page(1);
  });

  function refresh() {
    AccountsService.listTeamMembers(vm.team.id)
      .then(function(data) {
        vm.tableParams.settings({dataset: data});
      }, function(error) {
        MessageService.addErrorMessage('Error retrieving team members', error.data);
      });
    AccountsService.getTeam(vm.team.id)
      .then(function(data) {
        vm.team = data;
      }, function(error) {
        MessageService.addErrorMessage('Error retrieving team information', error.data);
      });
    AccountsService.list()
      .then(function(data) {
        vm.addUserTableParams.settings({dataset: data});
      }, function(error) {
        MessageService.addErrorMessage('Error retrieving users', error.data);
      });
  }
}

module.exports = TeamController;
