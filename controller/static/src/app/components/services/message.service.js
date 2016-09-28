'use strict';

MessageService.$inject = ['$rootScope'];
function MessageService($rootScope) {

  function addMessage (title, message, severity, log) {
    var newMsg = {
      title: title,
      message: message,
      severity: severity,
      log: log
    };

    // Create the messages array if it doesn't exist
    if(!$rootScope.messages) {
      $rootScope.messages = [];
    }

    // Don't add duplicate messages
    for(var i = 0; i < $rootScope.messages.length; i++) {
      if($rootScope.messages[i].title === newMsg.title &&
          $rootScope.messages[i].message === newMsg.message) {
        return;
      }
    }

    // Add the message to the array
    $rootScope.messages.push(newMsg);
  }

  function dismissMessage(msg) {
    for(var i = 0; i < $rootScope.messages.length; i++) {
      if($rootScope.messages[i].title === msg.title &&
          $rootScope.messages[i].message === msg.message) {
        $rootScope.messages.splice(i, 1);
      }
    }
  }

  return {
    addErrorMessage: function(title, message, log) {
      addMessage(title, message, 'error', log);
    },
    addSuccessMessage: function(title, message, log) {
      addMessage(title, message, 'success', log);
    },
    addInfoMessage: function(title, message, log) {
      addMessage(title, message, 'info', log);
    },
    addWarningMessage: function(title, message, log) {
      addMessage(title, message, 'warning', log);
    },
    dismissMessage: dismissMessage,
    clearMessages: function() {
      $rootScope.messages = [];
    }
  };
}

module.exports = MessageService;
