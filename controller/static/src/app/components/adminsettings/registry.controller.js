'use strict';

var $ = require('jquery');

RegistrySettingsController.$inject = ['RegistryService', 'MessageService', '$state'];
function RegistrySettingsController(RegistryService, MessageService, $state) {
  var vm = this;
  vm.registry = {
  };
  vm.caFile = '';
  vm.error = '';

  vm.updateRegistry = updateRegistry;
  vm.disableRegistry = disableRegistry;
  vm.showUpdateRegistryModal = showUpdateRegistryModal;
  vm.showDisableRegistryModal = showDisableRegistryModal;

  RegistryService.getRegistry()
    .then(function(data) {
      if(data.url) {
        vm.dtrEnabled = true;
      }
      vm.registry = data;
      vm.registry.name = 'Docker Trusted Registry';
      vm.registry.type = 'v2';
    },
    function(error) {
      vm.error = error.data;
    });

  function isFormValid() {
    return $('.ui.registry-config.form').form('validate form');
  }

  function showUpdateRegistryModal() {
    if(!isFormValid()) {
      return;
    }
    $('#update-registry-modal').modal('show');
  }

  function showDisableRegistryModal() {
    $('#disable-registry-modal').modal('show');
  }

  function hideUpdateRegistryModal() {
    $('#update-registry-modal').modal('hide');
  }

  function hideDisableRegistryModal() {
    $('#disable-registry-modal').modal('hide');
  }

  function refresh() {
    $state.go('dashboard.settings.registry', {}, { reload: true });
  }

  function disableRegistry() {
    RegistryService.updateRegistry({})
      .then(function(data) {
        refresh();
        MessageService.addSuccessMessage('Success', 'Disabled DTR integration');
        hideDisableRegistryModal();
      },
      function(error) {
        MessageService.addErrorMessage('Error disabling DTR integration', error.data);
      });
  }

  function updateRegistry() {
    vm.registry.ca_cert = vm.caFile;

    RegistryService.updateRegistry(vm.registry)
      .then(function(data) {
        refresh();

        MessageService.addSuccessMessage('Success', 'Updated DTR configuration');
        hideUpdateRegistryModal();
      },
      function(error) {
        MessageService.addErrorMessage('Error updating DTR configuration', error.data);
      });
  }
}

module.exports = RegistrySettingsController;
