'use strict';

TrustController.$inject = ['trustConfig', '$state', 'MessageService', 'TrustService'];
function TrustController(trustConfig, $state, MessageService, TrustService) {
  var vm = this;

  vm.trustConfig = trustConfig;
  vm.updateTrust = updateTrust;
  vm.require_content_trust = vm.trustConfig.require_content_trust_for_dtr || vm.trustConfig.require_content_trust_for_hub;

  vm.toggleRequireTrust = function() {
      if (vm.require_content_trust) {
          vm.trustConfig.require_content_trust_for_dtr = true;
          vm.trustConfig.require_content_trust_for_hub = true;
      } else {
          vm.trustConfig.require_content_trust_for_dtr = false;
          vm.trustConfig.require_content_trust_for_hub = false;
      }
  };

  function refresh() {
    $state.go('dashboard.settings.trust', {}, { reload: true });
  }

  function updateTrust() {
    TrustService.setTrustConfig(vm.trustConfig)
      .then(
        function(data) {
          refresh();

          MessageService.addSuccessMessage('Success', 'Updated content trust settings');
        },
        function(error) {
          MessageService.addErrorMessage('Error updating content trust settings', error.data);
        }
      );
  }
}

module.exports = TrustController;
