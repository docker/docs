'use strict';

var _ = require('lodash');
var $ = require('jquery');

ServicesController.$inject = ['nodes', 'TaskService', 'ServiceService', 'MessageService', 'NgTableParams', '$state', '$scope'];
function ServicesController(nodes, TaskService, ServiceService, MessageService, NgTableParams, $state, $scope) {
  var vm = this;
  vm.filter = '';
  vm.nodes = nodes;

  vm.tableParams = new NgTableParams({
    count: 25,
    filter: {
      $: vm.filter
    }
  });

  vm.removeService = removeService;
  vm.showRemoveServiceModal = showRemoveServiceModal;
  vm.inspectService = inspectService;

  // Global filtering
  $scope.$watch('vm.filter', function() {
    var filters = vm.tableParams.filter();
    filters.$ = vm.filter;
    vm.tableParams.filter(filters);
  });

  $scope.$watch('filter.$', function () {
      vm.tableParams.reload();
      vm.tableParams.page(1);
  });


  function showRemoveServiceModal(id) {
    $('#remove-service-modal')
      .modal({
        onApprove: function() {
          return vm.removeService(id);
        }
      })
      .modal('show');
  }

  function removeService(id) {
    ServiceService.remove(id)
      .then(function(response) {
        MessageService.addSuccessMessage('Successfully removed service');
        load();
      }, function(error) {
        MessageService.addErrorMessage('Error removing service', error.data);
      });
  }

  function inspectService(id) {
    $state.go('dashboard.resources.services.inspect', {id: id});
  }

  function load() {
    ServiceService.list()
      .then(function(services) {
        _.forEach(services, function(service) {
          if(service.Spec.Mode.Replicated) {
            service._Mode = 'Replicated';
          } else if(service.Spec.Mode.Global) {
            service._Mode = 'Global';
          }
        });

        // Turn our list of services into a map keyed by ID so that we can combine with our lists of tasks
        var servicesMap = _.keyBy(services, 'ID');

        TaskService.list()
          .then(function(tasks) {

            // Create status totals for each service
            for(var i = 0; i < tasks.length; i++) {
              if(!servicesMap[tasks[i].ServiceID]) {
                continue;
              }

              if(!servicesMap[tasks[i].ServiceID]._StatusSummary) {
                servicesMap[tasks[i].ServiceID]._StatusSummary = {
                    inactive: 0,
                    updating: 0,
                    active: 0,
                    errored: 0
                };
              }

              servicesMap[tasks[i].ServiceID]._StatusSummary[TaskService.stateSummaryMap[tasks[i].Status.State]]++;
            }

            vm.tableParams.settings({
              dataset: _.values(servicesMap)
            });
          }, function(error) {
            MessageService.addErrorMessage('Error retrieving services', error.data);
          });

      }, function(error) {
        MessageService.addErrorMessage('Error retrieving services', error.data);
      });

  }

  load();
}

module.exports = ServicesController;
