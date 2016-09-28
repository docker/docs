'use strict';

var $ = require('jquery');
var angular = require('angular');

var CodeMirror = require('codemirror');
require('../../../../../node_modules/codemirror/addon/mode/overlay.js');

DeployYmlController.$inject = ['ContainerService', 'MessageService', 'CAService', '$scope', '$timeout', '$state', 'store', '$http'];
function DeployYmlController(ContainerService, MessageService, CAService, $scope, $timeout, $state, store, $http) {
  /*global ORG*/
  var org = ORG || 'docker';
  /*global TAG*/
  var tag = TAG || 'latest';

  var vm = this;

  CodeMirror.defineMode('compose', function(config, parserConfig) {
    var compose = {
      token: function(stream, state) {
        var ch;
        if (
          stream.match('build') ||
          stream.match('dockerfile') ||
          stream.match('env_file')) {
          // consume until the end of line
          while ((ch = stream.next())) {} // eslint-disable-line no-empty

          return 'error';
        }

        // Consume until we find an 'error' token
        while (stream.next() &&
               !stream.match('build', false) &&
               !stream.match('dockerfile', false) &&
               !stream.match('env_file', false)) {} // eslint-disable-line no-empty
        return null;
      }
    };
    return CodeMirror.overlayMode(CodeMirror.getMode(config, parserConfig.backdrop || 'yaml'), compose);
  });

  vm.editorOpts = {
    indentWithTabs: false,
    lineWrapping: true,
    lineNumbers: true,
    theme: 'docker',
    mode: 'compose'
  };

  // Compose
  var composeImage = org + '/ucp-compose:' + tag;
  vm.projectName = '';
  vm.composeYml = '';
  vm.logs = '';
  vm.inProgress = false;

  vm.composeUp = composeUp;
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

  function composeUp() {
    if (!isFormValid()) {
      return;
    }

		vm.logs = '';
    vm.inProgress = true;

    var composeOpts = '';
    if(vm.composeVerbose) {
      composeOpts += '--verbose ';
    }

    CAService.getCACert().then(function(caCert) {
      var ucpComposeContainer = {
        Image: composeImage,
        AttachStdout: true,
        AttachStderr: true,
        Env: [
          'PROJECT_NAME=' + vm.projectName,
          'SESSION_TOKEN=' + store.get('sessionToken'),
          'CONTROLLER_HOST=' + window.location.hostname + ':' + (window.location.port || '443'),
          'DOCKER_COMPOSE_YML=' + btoa(vm.composeYml.replace(/\t/g, '  ')),
          'CONTROLLER_CA_CERT=' + btoa(caCert),
          'COMPOSE_OPTS=' + composeOpts
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
        vm.logs += 'Created ' + composeImage + ' compose container\n';

        $http
        .post('/containers/' + createData.Id + '/start', ucpComposeContainer)
        .success(function(startData) {
          vm.logs += 'Started compose container ' + createData.Id + '\n';

          // Wait for container to finish and then remove it
          $http
          .post('/containers/' + createData.Id + '/wait')
          .success(function(waitData) {
            if(waitData && waitData.StatusCode !== 0) {
              vm.logs += 'Deployment was unsuccessful (status code ' + waitData.StatusCode + ')\n';
              return;
            }

            vm.logs += 'Successfully deployed ' + vm.projectName + '\n';
          })
          .error(function(error) {
            vm.logs += 'Error deploying compose application: ' + error.data + '\n';
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
            vm.logs += 'Error getting logs from compose container: ' + error.data + '\n';
          });

        })
        .error(function(error) {
          vm.logs += 'Error starting compose container: ' + error.data + '\n';
          removeComposeContainer(createData.Id);
          vm.inProgress = false;
        });
      })
      .error(function(error) {
        vm.logs += 'Error creating compose container: ' + error + '\n';
        vm.inProgress = false;
      });

    }, function(error) {
      vm.logs = 'Unable to get CA certificate from controller\n';
      vm.inProgress = false;
    });
  }
}

module.exports = DeployYmlController;
