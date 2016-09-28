'use strict';

var nv = require('nvd3');
var d3 = require('d3');

DashboardController.$inject =
  ['DashboardService', 'ApplicationsService', 'TaskService', 'ServiceService', 'MessageService', '$scope'];
function DashboardController(
  DashboardService, ApplicationsService, TaskService, ServiceService, MessageService, $scope) {
  $scope.Math = Math;

  var vm = this;
  vm.refresh = refresh;
  vm.applicationCount = 0;
  vm.containerCount = 0;
  vm.imageCount = 0;
  vm.nodeCount = 0;
  vm.reservedCpuTotal = 0;
  vm.reservedCpuUsed = 0;
  vm.reservedCpu = 0;
  vm.reservedMemTotal = 0;
  vm.reservedMemUsed = 0.0;
  vm.strategy = 'Unknown';
  vm.strategyIconClass = 'help icon';
  vm.serviceToTasksMap = {};
  vm.numServices = 0;
  vm.runningServices = [];
  vm.partlyRunningServices = [];
  vm.preparingServices = [];
  vm.notRunningServices = [];

  refresh();

  function createDonutData(used, total) {
    return [
      {
        label: 'Free',
        value: total - used
      },
      {
        label: 'Used',
        value: used
      }
    ];
  }

  function createNvChart() {
    // FIXME: Can we get these charts extracted out into a directive
    // or some other way that we can re-use this, and only define colors
    // in a single place
    var chart = nv.models.pieChart()
      .x(function(d) {
        return d.label;
      })
      .y(function(d) {
        return d.value;
      })
      .showLabels(false)
      .showLegend(false)
      .donut(true)
      .donutRatio(0.35)
      .color(['#00CBCA', '#EF4A53', '#FFB463']);

    chart.tooltip.enabled();

    return chart;
  }

  function drawCpuDonut(used, total) {
    nv.addGraph(function() {
      var donutChart = createNvChart();

      // Insert text into the center of the donut
      function centerText() {
        return function() {
          var svg = d3.select('#cpu-pie-chart');

          var donut = svg.selectAll('g.nv-slice').filter(
            function (d, i) {
              return i === 0;
            }
          );

          donut.insert('text', 'g')
            .text('CPU')
            .attr('class', 'middle')
            .attr('text-anchor', 'middle')
            .attr('dy', '-.15em')
            .style('fill', '#000');

          donut.insert('text', 'g')
            .text(vm.reservedCpu + '%')
            .attr('class', 'middle')
            .attr('text-anchor', 'middle')
            .attr('dy', '1.0em')
            .style('fill', '#000');
        };
      }

      // Put the donut pie chart together
      d3.select('#cpu-pie-chart')
        .datum(createDonutData(vm.reservedCpuUsed, vm.reservedCpuTotal))
        .transition().duration(300)
        .call(donutChart)
        .call(centerText());

      d3.select('#cpu-pie-chart .nv-pieWrap')
        .attr('transform', 'scale(1.35) translate(-25, -30)');

      return donutChart;
    });
  }

  function drawMemoryDonut(used, total) {
    nv.addGraph(function() {
      var donutChart = createNvChart();

      // Insert text into the center of the donut
      function centerText() {
        return function() {
          var svg = d3.select('#memory-pie-chart');

          var donut = svg.selectAll('g.nv-slice').filter(
            function (d, i) {
              return i === 0;
            }
          );

          donut.insert('text', 'g')
            .text('Memory')
            .attr('class', 'middle')
            .attr('text-anchor', 'middle')
            .attr('dy', '-.15em')
            .style('fill', '#000');

          donut.insert('text', 'g')
            .text(vm.reservedMem + '%')
            .attr('class', 'middle')
            .attr('text-anchor', 'middle')
            .attr('dy', '1.0em')
            .style('fill', '#000');
        };
      }

      // Put the donut pie chart together
      d3.select('#memory-pie-chart')
        .datum(createDonutData(vm.reservedMemUsed, vm.reservedMemTotal))
        .transition().duration(300)
        .call(donutChart)
        .call(centerText());

      d3.select('#memory-pie-chart .nv-pieWrap')
        .attr('transform', 'scale(1.35) translate(-25, -30)');

      return donutChart;
    });
  }

  function refresh() {
    ApplicationsService.list()
      .then(function(data) {
        vm.applicationCount = data.length;
      }, function(error) {
        MessageService.addErrorMessage('Error getting applications', error.data);
      });

    DashboardService.info()
      .then(function(data) {
        vm.containerCount = data.Containers;
        vm.imageCount = data.Images;
        vm.systemStatus = DashboardService.parseSystemStatus(data.SystemStatus);
        vm.nodeCount = vm.systemStatus.nodeCount;
        vm.reservedCpu = vm.systemStatus.reservedCpu;
        vm.reservedCpuTotal = vm.systemStatus.reservedCpuTotal;
        vm.reservedCpuUsed = vm.systemStatus.reservedCpuUsed;
        vm.reservedMem = vm.systemStatus.reservedMem;
        vm.reservedMemTotal = vm.systemStatus.reservedMemTotal;
        vm.reservedMemUsed = vm.systemStatus.reservedMemUsed;
        vm.strategy = vm.systemStatus.strategy;
        vm.strategyIconClass = vm.systemStatus.strategyIconClass;

        drawCpuDonut(vm.reservedCpuUsed, vm.reservedCpuTotal);
        drawMemoryDonut(vm.reservedMemUsed, vm.reservedMemTotal);
      }, function(error) {
        MessageService.addErrorMessage('Error getting dashboard info', error.data);
      });

    DashboardService.nodeInfo()
      .then(function(data) {
        vm.nodes = data;
        if (Array.isArray(vm.nodes)) {
          vm.numNodes = vm.nodes.length;
          // Nodes that are ready
          vm.readyNodes = vm.nodes.filter(function(v) {
            return v.Status.State === 'ready';
          }).length;
          // Nodes that are down/UNREACHABLE
          vm.downNodes = vm.nodes.filter(function(v) {
            return v.Status.State === 'down';
          }).length;
          // Nodes that are MANAGERs
          vm.managerNodes = vm.nodes.filter(function(v) {
            return v.Spec.Role === 'manager';
          });
          // Nodes that are WORKERs
          vm.workerNodes = vm.nodes.filter(function(v) {
            return v.Spec.Role === 'worker';
          }).length;
          // Nodes that are MANAGERs and ready
          vm.readyManagerNodes = vm.nodes.filter(function(v) {
            return v.Spec.Role === 'manager' && v.Status.State === 'ready';
          }).length;
          // Nodes that are WORKERS and ready
          vm.readyWorkerNodes = vm.readyNodes - vm.readyManagerNodes;
        }
      }, function(error) {
        MessageService.addErrorMessage('Error getting swarm cluster nodes', error.data);
      });

    ServiceService.list()
      .then(function(data) {
        vm.services = data;

        if(Array.isArray(vm.services) && vm.services.length > 0) {
          vm.numServices = vm.services.length;
          TaskService.list()
            .then(function(tasksData) {
              vm.tasks = tasksData;
              vm.services.forEach(function(service) {
                vm.serviceToTasksMap[service.ID] = vm.tasks.filter(function(v) {
                  return v.ServiceID === service.ID;
                });
                var serviceTasks = vm.serviceToTasksMap[service.ID];
                if (serviceTasks) {
                  var runningTasks = serviceTasks.filter(function(t) {
                    return TaskService.stateSummaryMap[t.Status.State] === 'active';
                  });
                  if (runningTasks.length === serviceTasks.length) {
                    vm.runningServices.push(service);
                  } else if (runningTasks.length > 0 && runningTasks.length < serviceTasks.length) {
                    vm.partlyRunningServices.push(service);
                  } else {
                    var preparingTasks = serviceTasks.filter(function(t) {
                      return TaskService.stateSummaryMap[t.Status.State] === 'updating';
                    });
                    if (preparingTasks.length > 0) {
                      vm.preparingServices.push(service);
                    } else {
                      vm.notRunningServices.push(service);
                    }
                  }
                }

                // number of services
                vm.numRunningServices = vm.runningServices.length;
                vm.numPartlyRunningServices = vm.partlyRunningServices.length;
                vm.numPreparingServices = vm.preparingServices.length;
                vm.numNotRunningServices = vm.notRunningServices.length;
              });
            }, function(error) {
              MessageService.addErrorMessage('Error getting cluster-wide tasks information', error.data);
            });
        }
      }, function(error) {
        MessageService.addErrorMessage('Error getting cluster-wide services information', error.data);
    });

    DashboardService.containerInfo()
      .then(function(data) {
        vm.ctrs = data;
        if (Array.isArray(vm.ctrs)) {
          vm.numCtrs = vm.ctrs.length;
          vm.runningCtrs = vm.ctrs.filter(function(v) {
            return v.State === 'running';
          }).length;
          vm.stoppedCtrs = vm.ctrs.filter(function(v) {
            return v.State === 'exited';
          }).length;
          vm.pausedCtrs = vm.ctrs.filter(function(v) {
            return v.State === 'paused';
          }).length;
        }
      }, function(error) {
        MessageService.addErrorMessage('Error getting cluster-wide containers information', error.data);
      });
  }


}

module.exports = DashboardController;
