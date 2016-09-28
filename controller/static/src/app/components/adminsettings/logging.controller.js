'use strict';

var $ = require('jquery');

LoggingController.$inject = ['remoteLoggingConfig', '$state', 'MessageService', 'LoggingService'];
function LoggingController(remoteLoggingConfig, $state, MessageService, LoggingService) {
  var vm = this;

  vm.remoteLoggingConfig = remoteLoggingConfig;
  if(typeof vm.remoteLoggingConfig.host === 'undefined' || vm.remoteLoggingConfig.host === '') {
    vm.remoteLoggingEnabled = false;
  } else {
    vm.remoteLoggingEnabled = true;
  }

  vm.showDisableRemoteLoggingModal = showDisableRemoteLoggingModal;
  vm.showUpdateRemoteLoggingModal = showUpdateRemoteLoggingModal;

  vm.enableRemoteLogging = enableRemoteLogging;
  vm.disableRemoteLogging = disableRemoteLogging;

  function refresh() {
    $state.go('dashboard.settings.logging', {}, { reload: true });
  }

  function showUpdateRemoteLoggingModal() {
    $('#update-remote-logging-modal').modal('show');
  }

  function showDisableRemoteLoggingModal() {
    $('#disable-remote-logging-modal').modal('show');
  }

  function hideUpdateRemoteLoggingModal() {
    $('#update-remote-logging-modal').modal('hide');
  }

  function hideDisableRemoteLoggingModal() {
    $('#disable-remote-logging-modal').modal('hide');
  }

  function enableRemoteLogging() {
    LoggingService.setRemoteLoggingConfig(vm.remoteLoggingConfig)
      .then(
        function(data) {
          refresh();

          MessageService.addSuccessMessage('Success', 'Updated logging configuration');
          hideUpdateRemoteLoggingModal();
        }, function(error) {
          MessageService.addErrorMessage('Error updating logging configuration', error.data);
        });
    hideDisableRemoteLoggingModal();
  }

  function disableRemoteLogging() {
    LoggingService.disableRemoteLogging()
      .then(
        function(data) {
          refresh();
          MessageService.addSuccessMessage('Success', 'Disabled remote logging');
          hideDisableRemoteLoggingModal();
        },
        function(error) {
          MessageService.addErrorMessage('Error disabling remote logging', error.data);
        }
      );
  }
}

module.exports = LoggingController;
