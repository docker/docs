'use strict';

var Terminal = require('term.js');
var $ = require('jquery');

ExecController.$inject = ['$http', '$state', '$stateParams', 'ContainerService', 'MessageService'];
function ExecController($http, $state, $stateParams, ContainerService, MessageService) {
  $state.current.ncyBreadcrumb.parent = 'dashboard.inspect.details({id: \'' + $stateParams.id + '\'})';

  var vm = this;
  vm.id = $stateParams.id;
  vm.addr = '';
  vm.command = 'bash';
  vm.connected = false;

  vm.connect = connect;
  vm.disconnect = disconnect;

  var term;
  var websocket;

  function connect() {
    var termWidth = Math.round($(window).width() / 7.5);
    var termHeight = 30;
    var cmd = vm.command.replace(' ', ',');

    var url = window.location.href;
    var urlparts = url.split('/');
    var scheme = urlparts[0];
    var wsScheme = 'ws';

    if (scheme === 'https:') {
      wsScheme = 'wss';
    }

    // we make a request for a console session token; this is used
    // as authentication to make sure the user has console access
    // for this exec session
    $http.post('/api/consolesession/' + vm.id)
      .success(function(data, status, headers, config) {
        vm.connected = true;
        vm.token = data.token;
        vm.addr = wsScheme + '://' + window.location.hostname + ':' + window.location.port + '/exec?id=' + vm.id + '&cmd=' + cmd + '&h=' + termHeight + '&w=' + termWidth + '&token=' + vm.token;

        if (typeof term !== 'undefined' && term !== null) {
          term.destroy();
        }

        websocket = new WebSocket(vm.addr);

        websocket.onopen = function(evt) {
          term = new Terminal({
            cols: termWidth,
            rows: termHeight,
            screenKeys: true,
            useStyle: true,
            cursorBlink: true
          });
          term.on('data', function(termData) {
            websocket.send(termData);
          });
          term.open(document.getElementById('container-terminal'));
          websocket.onmessage = function(msg) {
            term.write(msg.data);
          };
          websocket.onclose = function(close) {
            term.write('Session terminated');
            term.destroy();
            vm.connected = false;
          };
          websocket.onerror = function(error) {
            MessageService.addErrorMessage('Error occurred during exec session', error.data);
            vm.connected = false;
          };
        };
      })
      .error(function(error) {
        MessageService.addErrorMessage('Error attempting to open exec session', error.data);
      });
  }

  function disconnect() {
    if (websocket !== null) {
      websocket.close();
    }

    if (term) {
      term.destroy();
    }
    vm.connected = false;
  }
}

module.exports = ExecController;
