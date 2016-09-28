'use strict';

var oboe = require('oboe');
var $ = require('jquery');

ImagesController.$inject = ['ImagesService', 'MessageService', '$state', '$timeout', '$scope', 'LoaderService', 'NgTableParams', 'AuthService'];
function ImagesController(ImagesService, MessageService, $state, $timeout, $scope, LoaderService, NgTableParams, AuthService) {
  var vm = this;
  vm.user = AuthService.getCurrentUser();
  vm.pulling = false;
  vm.selectedImage = null;
  vm.pullImageName = '';
  vm.filter = '';
  vm.tableParams = new NgTableParams({
    count: 25
  }, {});

  vm.refresh = refresh;
  vm.removeImage = removeImage;
  vm.pullImage = pullImage;
  vm.showRemoveImageDialog = showRemoveImageDialog;
  vm.showPullImageDialog = showPullImageDialog;

  refresh();

  function showRemoveImageDialog(image) {
    vm.selectedImage = image;
    $('#remove-modal').modal('show');
  }

  function hidePullImageDialog() {
    vm.refresh();
    LoaderService.stop();
  }

  function showPullImageDialog(image) {
    $('#pull-modal')
      .modal({
        onApprove: function() {
          return vm.pullImage();
        },
        onDeny: function() {
          return hidePullImageDialog();
        }
      })
      .modal('show');
  }

  // Global filtering
  $scope.$watch('vm.filter', function() {
    vm.tableParams.filter({ $: vm.filter });
  });

  $scope.$watch('filter.$', function () {
    vm.tableParams.reload();
    vm.tableParams.page(1);
  });

  function refresh() {
    vm.pullImageName = '';
    vm.pullImageError = '';
    ImagesService.list()
      .then(function(data) {
        vm.tableParams.settings({dataset: data});
      }, function(error) {
        MessageService.addErrorMessage('Error retrieving images', error.data);
      });
  }

  function removeImage() {
    ImagesService.remove(vm.selectedImage)
      .then(function(data) {
        MessageService.addSuccessMessage('Success', 'Removed image ' + vm.selectedImage.RepoTags[0]);
        vm.refresh();
      }, function(data) {
        MessageService.addErrorMessage('Error removing image', data.data);
      });
  }

  function pullImage() {
    if (!$('#pull-modal .ui.form').form('validate form')) {
      return false;
    }

    LoaderService.setMessage('Pulling Image: ' + vm.pullImageName);
    LoaderService.start();

    var headers = {
      'Authorization': 'Bearer ' + vm.user.sessionToken
    };

    if(vm.registryUsername && vm.registryPassword) {
      headers['X-Registry-Auth'] = btoa(JSON.stringify({
        username: vm.registryUsername,
        password: vm.registryPassword
      }));
    }

    oboe({
      url: '/images/create?fromImage=' + vm.pullImageName,
      method: 'POST',
      withCredentials: true,
      headers: headers
    })
    .done(function(node) {
      if(node.error) {
        vm.pullImageError = node.error;
        LoaderService.stop();
        showPullImageDialog();
        $timeout(function() {
          $scope.$apply();
        });
      } else if(node.status && node.status.indexOf(':') > -1) {
        // We expect two nodes, e.g.
        // 1) Pulling busybox...
        // 2) Pulling busybox... : downloaded
        if(node.status.indexOf('downloaded') > -1) {
          MessageService.addSuccessMessage('Success', 'Pulled image ' + vm.pullImageName);

          // If we get a downloaded message, we're done
          vm.refresh();
          LoaderService.stop();
          hidePullImageDialog();
        }
      }

    })
    .fail(function(error) {
      setTimeout(function() {
        if(error.body && error.statusCode) {
          vm.pullImageError = error.body;
          showPullImageDialog();
        }
        LoaderService.stop();
        $timeout(function() {
          $scope.$apply();
        });
      });
    });

  return true;
  }
}

module.exports = ImagesController;
