'use strict';

var $ = require('jquery');

UserController.$inject = ['account', 'permissions', 'AccountsService', 'MessageService', '$scope', '$state', '$http', 'AuthService', 'NgTableParams'];
function UserController(account, permissions, AccountsService, MessageService, $scope, $state, $http, AuthService, NgTableParams) {
  $state.current.ncyBreadcrumb.label = 'Profile: ' + AuthService.getCurrentUser().username;

  var vm = this;
  vm.account = account;
  vm.tableParams = new NgTableParams({}, {
    dataset: permissions
  });

  vm.label = '';
  vm.publicKey = '';

  var passwordForm = $('.ui.password.form');
  vm.newLabel = '';
  vm.changeRequest = {
    old_password: '',
    new_password: ''
  };

  vm.changePassword = changePassword;
  vm.createAndDownloadBundle = createAndDownloadBundle;
  vm.showRenameModal = showRenameModal;
  vm.showRemoveModal = showRemoveModal;
  vm.showAddPublicKeyModal = showAddPublicKeyModal;
  vm.renameBundle = renameBundle;
  vm.removeBundle = removeBundle;
  vm.addPublicKey = addPublicKey;
  vm.roleDescription = roleDescription;

  // Apply jQuery to dropdowns in table once ngRepeat has finished rendering
  $scope.$on('ngRepeatFinished', function() {
    $('.ui.dropdown').dropdown();
  });

  function roleDescription(roleId) {
    if(roleId === 0) {
      return 'No Access';
    } else if(roleId === 1) {
      return 'View Only';
    } else if(roleId === 2) {
      return 'Restricted Control';
    } else if(roleId === 3) {
      return 'Full Control';
    }
    return 'Unknown';
  }

  function refresh() {
    $state.go('dashboard.user', {}, { reload: true });
  }

  function addPublicKey() {
    if(typeof vm.account.public_keys === 'undefined') {
      vm.account.public_keys = [];
    }

    // Clone the account and update it
    var updatedAccount = $.extend(true, {}, vm.account);
    updatedAccount.public_keys.push({label: vm.label, public_key: vm.publicKey});

    AccountsService.update(updatedAccount)
      .then(
        function(data) {
          refresh();

          MessageService.addSuccessMessage('Success', 'Added public key');
        },
        function(error) {
          MessageService.addErrorMessage('Error adding public key', error.data);
        }
      );
  }

  function renameBundle() {
    for(var i = 0; i < vm.account.public_keys.length; i++) {
      if(vm.account.public_keys[i].public_key === vm.publicKeyAction) {
        vm.account.public_keys[i].label = vm.newLabel;
      }
    }
    AccountsService.update(vm.account)
      .then(
        function(data) {
          refresh();

          MessageService.addSuccessMessage('Success', 'Renamed bundle');
          hideRenameModal();
        },
        function(error) {
          MessageService.addErrorMessage('Error renaming bundle', error.data);
        }
      );
  }

  function removeBundle() {
    for(var i = 0; i < vm.account.public_keys.length; i++) {
      if(vm.account.public_keys[i].public_key === vm.publicKeyAction) {
        vm.account.public_keys.splice(i, 1);
      }
    }
    AccountsService.update(vm.account)
      .then(
        function(data) {
          refresh();

          MessageService.addSuccessMessage('Success', 'Removed bundle');
          hideRemoveModal();
        },
        function(error) {
          MessageService.addErrorMessage('Error removing bundle', error.data);
        }
      );
  }

  function showRenameModal(publicKey) {
    vm.publicKeyAction = publicKey;
    for(var i = 0; i < vm.account.public_keys.length; i++) {
      if(vm.account.public_keys[i].public_key === vm.publicKeyAction) {
        vm.newLabel = vm.account.public_keys[i].label;
      }
    }
    $('#rename-modal').modal('show');
  }

  function showRemoveModal(publicKey) {
    vm.publicKeyAction = publicKey;
    $('#remove-modal').modal('show');
  }

  function hideRemoveModal() {
    $('#remove-modal').modal('hide');
  }

  function hideRenameModal() {
    $('#rename-modal').modal('hide');
  }

  function showAddPublicKeyModal() {
    $('#add-public-key-modal').modal('show');
  }

  function createAndDownloadBundle() {
    $http
      .get('/api/clientbundle', { responseType: 'arraybuffer' })
      .then(function(response) {
        var file = new Blob([ response.data ], {
          type: 'application/zip'
        });
        var filename = 'orca-bundle.zip';
        var rx = response.headers('Content-Disposition').match(/inline; filename='(.*?)'/);
        if(rx.length > 0) {
          filename = rx[1];
        }
        //trick to download store a file having its URL
        var fileURL = URL.createObjectURL(file);
        var a = document.createElement('a');
        a.href = fileURL;
        a.target = '_blank';
        a.download = filename;
        document.body.appendChild(a);
        a.click();

        refresh();
      },
      function(error) {
        MessageService.addErrorMessage('Error creating bundle', error.statusText);
      });
  }

  function isFormValid() {
    return passwordForm.form('validate form');
  }

  function changePassword() {
    if(!isFormValid()) {
      return;
    }

    $http({
      method: 'POST',
      url: '/account/changepassword',
      data: vm.changeRequest,
      ignore401: true
    }).then(function(success) {
      $state.go('logout');
    }, function(error) {
      MessageService.addErrorMessage('Error changing password', error.data);
    });
  }
}

module.exports = UserController;
