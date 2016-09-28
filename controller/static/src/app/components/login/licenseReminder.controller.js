'use strict';

var $ = require('jquery');

LicenseReminderController.$inject = ['LicensingService', 'MessageService', '$state'];
function LicenseReminderController(LicensingService, MessageService, $state) {
  var vm = this;
  /*global REQUIRE_LICENSE*/
  vm.requireLicense = (REQUIRE_LICENSE === '1' || REQUIRE_LICENSE === 'true' ? true : false);
  vm.error = '';

  vm.uploadLicense = uploadLicense;
  vm.uploadFileChange = uploadFileChange;

  function uploadFileChange(e) {
    var licenseConfig = {};
    var uploadedFile = e.files[0];
    var reader = new FileReader();
    reader.readAsText(uploadedFile, 'UTF-8');
    reader.onload = function (evt) {
      var uploadedFileContent = evt.target.result;
      try {
        licenseConfig = JSON.parse(uploadedFileContent);
        var licenseRequest = {
          auto_refresh: true,
          license_config: licenseConfig
        };
        LicensingService.updateLicense(licenseRequest)
          .then(function(data) {
            $state.go('dashboard.main');
            MessageService.addSuccessMessage('Success', 'UCP license uploaded');
          },
          function(error) {
            vm.error = 'Error updating license: ' + error.data;
          });
      } catch(jsonParseError) {
        vm.error = 'Could not parse license file: ' + jsonParseError;
      }
    };
    reader.onerror = function () {
      vm.error = 'Error reading file.';
    };
  }

  function uploadLicense() {
    $('.license-file-input').click();
  }

}

module.exports = LicenseReminderController;
