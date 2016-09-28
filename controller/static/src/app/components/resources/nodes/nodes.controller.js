'use strict';

var $ = require('jquery');
var _ = require('lodash');

NodeController.$inject = ['info', 'AuthService', 'ServiceService', 'TaskService', 'NodeService', 'MessageService', '$q', '$scope', '$state', '$timeout'];
function NodeController(info, AuthService, ServiceService, TaskService, NodeService, MessageService, $q, $scope, $state, $timeout) {
  var websocket;
  var agentServiceName = 'ucp-agent';

  var vm = this;
  vm.info = info;
  vm.user = AuthService.getCurrentUser();
  vm.filter = '';
  vm.selectedNode = {};
  vm.addNode = {};

  vm.inspectNode = inspectNode;
  vm.removeNode = removeNode;
  vm.activate = activate;
  vm.pause = pause;
  vm.drain = drain;
  vm.promote = promote;
  vm.demote = demote;
  vm.showAddNodeModal = showAddNodeModal;
  vm.showRemoveNodeModal = showRemoveNodeModal;

  $scope.$on('ngRepeatFinished', function() {
    $('.tooltip.icon').popup({
      inline: true,
      hoverable: true,
      position: 'bottom center',
      delay: {
        show: 150,
        hide: 150
      }
    });
  });

  // Reset the clipboard copy when something changes
  $scope.$watch('vm.addNode', function() {
    if(vm.addNode && vm.addNode.useCustomListenAddr === false) {
      vm.addNode.listenAddr = '';
    }
    vm.copied = false;
  }, true);

  $scope.$on('$destroy', function() {
    if(websocket) {
      websocket.close();
    }
  });

  function showAddNodeModal() {
    $('#add-node-modal').modal('show');
  }

  function inspectNode(id) {
    $state.go('dashboard.resources.nodes.inspect', {id: id});
  }

  function createBeachheadContainerIDMapping() {
    ServiceService.list()
      .then(function(services) {
        var servicesMap = _.keyBy(services, 'Spec.Name');
        if(!servicesMap[agentServiceName]) {
          MessageService.addErrorMessage('Could not find the UCP agent service');
          return;
        }

        var serviceId = servicesMap[agentServiceName].ID;

        TaskService.list()
          .then(function(tasks) {
            var filteredTasks = _.filter(tasks, function(t) {
              return t.ServiceID === serviceId;
            });

            vm.beachheadTasks = _.keyBy(filteredTasks, 'NodeID');
            _.forEach(filteredTasks, function(t) {
              if(vm.beachheadTasks[t.NodeID] && new Date(vm.beachheadTasks[t.NodeID].CreatedAt) > new Date(t.CreatedAt)) {
                return;
              }

              vm.beachheadTasks[t.NodeID] = t;
            });
          });
      });
  }

  function showRemoveNodeModal(node) {
    vm.selectedNode = node;
    $('#remove-node-modal').modal('show');
  }

  function activate(id) {
    NodeService.activate(id)
      .then(function(data) {
        load();
        MessageService.addSuccessMessage('Successfully activated node');
      }, function(error) {
        MessageService.addErrorMessage('Error activating node', error.data.message);
      });
  }

  function pause(id) {
    NodeService.pause(id)
      .then(function(data) {
        load();
        MessageService.addSuccessMessage('Successfully paused node');
      }, function(error) {
        MessageService.addErrorMessage('Error pausing node', error.data.message || error.data);
      });
  }

  function drain(id) {
    NodeService.drain(id)
      .then(function(data) {
        load();
        MessageService.addSuccessMessage('Successfully set node to drain');
      }, function(error) {
        MessageService.addErrorMessage('Error setting node to drain', error.data.message);
      });
  }

  function promote(id) {
    NodeService.promote(id)
      .then(function(data) {
        load();
        MessageService.addSuccessMessage('Successfully promoted node');
      }, function(error) {
        MessageService.addErrorMessage('Error promoting node', error.data.message);
      });
  }

  function demote(id) {
    NodeService.demote(id)
      .then(function(data) {
        load();
        MessageService.addSuccessMessage('Successfully demoted node');
      }, function(error) {
        MessageService.addErrorMessage('Error demoting node', error.data.message);
      });
  }

  function removeNode(id) {
    NodeService.remove(id)
      .then(function(data) {
        load();
        MessageService.addSuccessMessage('Node removed');
      }, function(error) {
        MessageService.addErrorMessage('Error removing', error.data.message);
      });
  }

  function load() {
    $q.all([
      NodeService.listClassic(),
      NodeService.list()
    ]).then(function(data) {
      var i = 1;
      _.forEach(data[1], function(n) {
        if(!n.Description.Hostname) {
          n.Description.Hostname = 'Unknown (' + i++ + ')';
        }
      });
      var classic = _.keyBy(data[0], 'name');
      var nodes = _.keyBy(data[1], 'Description.Hostname');
      vm.nodes = _.values(_.merge(classic, nodes));
    }, function(error) {
      MessageService.addErrorMessage('Error retrieving nodes', error.data);
    });

    createBeachheadContainerIDMapping();
  }

  load();
}

module.exports = NodeController;
