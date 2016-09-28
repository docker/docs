'use strict';

SchedulerController.$inject = ['schedulerConfig', '$state', 'MessageService', 'SchedulerService'];
function SchedulerController(schedulerConfig, $state, MessageService, SchedulerService) {
  var vm = this;

  vm.schedulerConfig = schedulerConfig;
  vm.updateScheduler = updateScheduler;

  function refresh() {
    $state.go('dashboard.settings.scheduler', {}, { reload: true });
  }

  function updateScheduler() {
    SchedulerService.setSchedulerConfig(vm.schedulerConfig)
      .then(
        function(data) {
          refresh();

          MessageService.addSuccessMessage('Success', 'Updated scheduler settings');
        },
        function(error) {
          MessageService.addErrorMessage('Error updating scheduler settings', error.data);
        }
      );
  }
}

module.exports = SchedulerController;
