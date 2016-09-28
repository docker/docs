'use strict';

TrackingController.$inject = ['trackingConfig', '$state', 'MessageService', 'TrackingService'];
function TrackingController(trackingConfig, $state, MessageService, TrackingService) {
  var vm = this;

  vm.trackingConfig = trackingConfig;
  vm.updateTracking = updateTracking;

  function refresh() {
    $state.go('dashboard.settings.tracking', {}, { reload: true });
  }

  function updateTracking() {
    TrackingService.setTrackingConfig(vm.trackingConfig)
      .then(
        function(data) {
          refresh();

          MessageService.addSuccessMessage('Success', 'Updated usage reporting settings');
        },
        function(error) {
          MessageService.addErrorMessage('Error updating usage reporting', error.data);
        }
      );
  }
}

module.exports = TrackingController;
