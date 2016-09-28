'use strict';

var $ = require('jquery');
var _ = require('lodash');
var angular = require('angular');

DeployStackController.$inject = ['ContainerService', 'MessageService', 'CAService', '$scope', '$timeout', '$state', 'store', '$http'];
function DeployStackController(ContainerService, MessageService, CAService, $scope, $timeout, $state, store, $http) {
  /*global ORG*/
  var org = ORG || 'docker';
  /*global TAG*/
  var tag = TAG || 'latest';
  $scope._ = _;
  var vm = this;

  // Compose
  var composeImage = org + '/ucp-compose:' + tag;
  vm.projectName = '';

  vm.bundleFile = null;

  vm.logs = '';
  vm.inProgress = false;

  vm.createStack = createStack;
  vm.composeBack = composeBack;
  vm.removeComposeContainer = removeComposeContainer;

  function composeBack() {
    vm.logs = '';
  }

  function isFormValid() {
    return $('#deploy').form('validate form');
  }

  function removeComposeContainer(id) {
    ContainerService.remove(id)
      .then(function(data) {},
      function(error) {
        vm.logs += 'Error removing compose container: ' + error.data + '\n';
      });
  }
    $scope.$watch('vm.file', function(value) {
      if (value) {
        var r = new FileReader();
        r.onload = function(e) {
          vm.bundleFile = e.target.result;
          try {
            vm.parsedBundle = JSON.parse(e.target.result);
          } catch(error) {
            MessageService.addErrorMessage('Invalid bundle file', error);
          }
        };
        r.readAsText(value);
      }
    });

  function createStack() {
    if (!isFormValid()) {
      return;
    }

		vm.logs = '';
    vm.inProgress = true;
    vm.statusCode = null;

    CAService.getCACert().then(function(caCert) {
      var ucpComposeContainer = {
        Image: composeImage,
        AttachStdout: true,
        AttachStderr: true,
        Entrypoint: 'bundle-wrapper.sh',
        Env: [
          'NAMESPACE=' + vm.projectName,
          'DOCKER_STACK_BUNDLE=' + btoa(vm.bundleFile),
          'SESSION_TOKEN=' + store.get('sessionToken'),
          'CONTROLLER_HOST=' + window.location.hostname + ':' + (window.location.port || '443'),
          'CONTROLLER_CA_CERT=' + btoa(caCert)
        ],
				Config: {},
				HostConfig: {},
				NetworkConfig: {}
      };

      var wsScheme = 'ws';
      if (window.location.protocol === 'https:') {
        wsScheme = 'wss';
      }

      $http
      .post('/containers/create', ucpComposeContainer)
      .success(function(createData, createStatus, createHeaders, createConfig) {

        $http
        .post('/containers/' + createData.Id + '/start', ucpComposeContainer)
        .success(function(startData) {

          // Wait for container to finish and then remove it
          $http
          .post('/containers/' + createData.Id + '/wait')
          .success(function(waitData) {
            vm.statusCode = waitData.StatusCode;
            if(waitData && waitData.StatusCode !== 0) {
              MessageService.addErrorMessage('Deployment failed', 'Deployment was unsuccessful (status code ' + waitData.StatusCode + ')');
              return;
            }
            $state.go('dashboard.resources.services');
            MessageService.addSuccessMessage('Deployed ' + vm.projectName);
          })
          .error(function(error) {
            MessageService.addErrorMessage('Deployment failed', error.data);
          }).finally(function() {
            vm.inProgress = false;
            removeComposeContainer(createData.Id);
          });

          // Start streaming the logs from the compose container
          $http.post('/api/containerlogs/' + createData.Id)
          .success(function(logsession) {
            vm.addr = wsScheme + '://' + window.location.hostname + ':' + window.location.port + '/containerlogs?token=' + logsession.token + '&id=' + createData.Id;

            var websocket = new WebSocket(vm.addr);

            websocket.onopen = function(evt) {
              websocket.onmessage = function(msg) {
                vm.logs += msg.data;
                $('#compose-modal').modal('refresh');
                $timeout(function() {
                  $scope.$apply();
                });
              };
              websocket.onerror = function(error) {
                vm.logs += 'WebSocket error on getting logs: ' + error.data + '\n';
              };
            };
          })
          .error(function(error) {
            vm.logs += 'Error getting logs: ' + error.data + '\n';
          });

        })
        .error(function(error) {
          vm.logs += 'Error starting container: ' + error.data + '\n';
          removeComposeContainer(createData.Id);
          vm.inProgress = false;
        });
      })
      .error(function(error) {
        vm.logs += 'Error creating container: ' + error + '\n';
        vm.inProgress = false;
      });

    }, function(error) {
      vm.logs = 'Unable to get CA certificate from controller\n';
      vm.inProgress = false;
    });
  }
}

module.exports = DeployStackController;
