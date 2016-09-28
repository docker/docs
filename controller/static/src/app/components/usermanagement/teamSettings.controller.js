'use strict';

var $ = require('jquery');

TeamSettingsController.$inject = ['team', '$state', 'AccountsService', 'MessageService'];
function TeamSettingsController(team, $state, AccountsService, MessageService) {
  var vm = this;
  vm.team = team;

  vm.newTeamName = '';

  vm.showRenameTeamModal = showRenameTeamModal;
  vm.renameTeam = renameTeam;
  vm.showDeleteTeamModal = showDeleteTeamModal;
  vm.deleteTeam = deleteTeam;

  function showRenameTeamModal() {
    $('#rename-team-modal').modal('show');
  }

  function refresh() {
    $state.go('dashboard.accounts.teamSettings', { teamId: vm.team.id }, { reload: true });
  }

  function renameTeam() {
    var updatedTeam = vm.team;
    updatedTeam.name = vm.newTeamName;

    AccountsService.updateTeam(updatedTeam)
      .then(function(data) {
        refresh();

        MessageService.addSuccessMessage('Success', 'Renamed team to ' + vm.newTeamName);
      }, function(error) {
        MessageService.addErrorMessage('Error renaming team', error.data);
      });
  }

  function showDeleteTeamModal() {
    $('#delete-team-modal').modal('show');
  }

  function deleteTeam() {
    AccountsService.deleteTeam(vm.team.id)
      .then(function(data) {
        // Upon removing a team, redirect back to 'All Users' page.
        $state.go('dashboard.accounts.users', {}, { reload: true });

        MessageService.addSuccessMessage('Success', 'Removed team ' + vm.team.name);
      }, function(error) {
        MessageService.addErrorMessage('Error removing team', error.data);
      });
  }

}

module.exports = TeamSettingsController;
