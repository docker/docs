'use strict';

CertsController.$inject = ['CertsService', 'MessageService', '$state'];
function CertsController(CertsService, MessageService, $state) {
  var vm = this;
  vm.updateRequest = {
    ca: '',
    key: '',
    cert: ''
  };

  vm.uploadCerts = uploadCerts;

  function uploadCerts() {
    CertsService.updateCerts(vm.updateRequest)
      .then(function(data) {
        MessageService.addSuccessMessage('Success', 'Updated certs');
      },
      function(error) {
        MessageService.addErrorMessage('Error updating certs', error.data);
      });
  }
}

module.exports = CertsController;
