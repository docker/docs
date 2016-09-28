'use strict';

LicensingController.$inject = ['LicensingService', 'MessageService', '$state'];
function LicensingController(LicensingService, MessageService, $state) {
  var vm = this;
  vm.license = {};
  vm.licenseFile = '';

  vm.uploadLicense = uploadLicense;
  vm.setAutoRefresh = setAutoRefresh;

  LicensingService.getLicense()
    .then(function(data) {
      vm.license = data;
    },
    function(error) {
      MessageService.addErrorMessage('Error retrieving license information', error.data);
    });

  function refresh() {
    $state.go('dashboard.settings.licensing', {}, { reload: true });
  }

  function setAutoRefresh() {
    LicensingService.updateLicense(vm.license)
      .then(function(data) {
        refresh();

        MessageService.addSuccessMessage('Success', 'Updated license auto-refresh setting');
      },
      function(error) {
        MessageService.addErrorMessage('Error updating license auto-refresh', error.data);
      });
  }

  function uploadLicense() {
    var licenseConfig = {};
    try {
      licenseConfig = JSON.parse(vm.licenseFile);
    } catch(e) {
      MessageService.addErrorMessage('Could not parse license file', e);
    }

    var licenseRequest = {
      auto_refresh: vm.license.auto_refresh,
      license_config: licenseConfig
    };

    LicensingService.updateLicense(licenseRequest)
      .then(function(data) {
        refresh();

        MessageService.addSuccessMessage('Success', 'Updated license');
      },
      function(error) {
        MessageService.addErrorMessage('Error updating license', error.data);
      });
  }

}

module.exports = LicensingController;
