'use strict';

var $ = require('jquery');

AccountsBaseController.$inject = ['teams', '$scope', '$state', '$timeout', 'AccountsService', 'MessageService'];
function AccountsBaseController(teams, $scope, $state, $timeout, AccountsService, MessageService) {
  var vm = this;
  vm.teams = teams.sort(sortTeams);

  vm.selectedTeam = 'All Users';
  vm.selectedTeam = $state.params.teamId || 'All Users';
  $timeout(function() {
    $('.ui.dropdown').dropdown('set selected', vm.selectedTeam);
  });

  vm.refreshTeams = refreshTeams;
  vm.showCreateTeam = showCreateTeam;
  vm.newTeamType = 'Managed';
  vm.newTeamRequest = {};
  vm.addTeam = addTeam;

  function addTeam() {
    if (!$('.ui.new-team.form').form('validate form')) {
      return false;
    }
    AccountsService.listTeams()
      .then(function(data) {
        // Check if this team name already exists
        for(var i = 0; i < data.length; i++) {
          if(data[i].name === vm.newTeamRequest.name) {
            MessageService.addErrorMessage('Team Already Exists', 'A team with the name \'' + vm.newTeamRequest.name + '\' already exists');
            return;
          }
        }

        AccountsService.createTeam(vm.newTeamRequest)
          .then(function(response) {
            // Upon successful team creation, redirect to the page of the newly created team
            $state.go('dashboard.accounts.team', { teamId: response.data.id }, { reload: true });

            MessageService.addSuccessMessage('Success', 'Created team ' + vm.newTeamRequest.name);
          }, function(error) {
            MessageService.addErrorMessage('Error creating team', error.data);
          });
      }, function(error) {
        MessageService.addErrorMessage('Error creating team', error.data);
      });

    return true;
  }

  function sortTeams(a, b) {
    var textA = a.name.toUpperCase();
    var textB = b.name.toUpperCase();
    return (textA < textB) ? -1 : (textA > textB) ? 1 : 0;
  }

  function refreshTeams() {
    AccountsService.listTeams()
      .then(function(data) {
        vm.teams = data.sort(sortTeams);
      }, function(error) {
        MessageService.addErrorMessage('Error retrieving teams', error.data);
      });
  }

  function showCreateTeam() {
    $('#create-team-modal')
      .modal({
        onApprove: function() {
          return vm.addTeam();
        }
      })
      .modal('show');
  }

  $scope.$watch('vm.selectedTeam', function(newValue, oldValue) {
    if(newValue && newValue !== 'All Users') {
      vm.teamView = true;
      if(newValue !== oldValue) {
        $state.go('dashboard.accounts.team', { teamId: newValue });
      }
    } else {
      vm.teamView = false;
      if(newValue !== oldValue) {
        $state.go('dashboard.accounts.users');
      }
    }
  });
}

module.exports = AccountsBaseController;
