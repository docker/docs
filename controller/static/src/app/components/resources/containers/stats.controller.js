'use strict';

var $ = require('jquery');
var oboe = require('oboe');
var nv = require('nvd3');
var d3 = require('d3');

ContainerStatsController.$inject = ['$stateParams', '$scope', '$state', 'MessageService', 'store'];
function ContainerStatsController($stateParams, $scope, $state, MessageService, store) {
  $state.current.ncyBreadcrumb.parent = 'dashboard.inspect.details({id: \'' + $stateParams.id + '\'})';

  var vm = this;
  var graphCpuStats;
  var graphMemoryStats;
  var graphNetStats;
  vm.refreshInterval = 2;
  vm.id = $stateParams.id;

  vm.cpuStats = [{
    key: 'CPU',
    values: []
  }];

  vm.memoryStats = [{
    key: 'Memory',
    color: '#ffa500',
    values: []
  }];

  vm.netStats = [
    {
      key: 'Network (rx)',
      values: []
    },
    {
      key: 'Network (tx)',
      values: []
    }
  ];

  var previousCpuUsage = 0;
  var previousSystemCpuUsage = 0;
  function addCpuUsage(date, systemUsage, usage, cpuCores) {
    if(previousCpuUsage === 0 || previousSystemCpuUsage === 0) {
      previousCpuUsage = usage;
      previousSystemCpuUsage = systemUsage;
      return;
    }

    var usageSample = usage - previousCpuUsage;
    previousCpuUsage = usage;

    var systemUsageSample = systemUsage - previousSystemCpuUsage;
    previousSystemCpuUsage = systemUsage;

    var cpuPercent = 0.0;
    if(usageSample > 0.0 && systemUsageSample > 0.0) {
      cpuPercent = (usageSample / systemUsageSample) * cpuCores * 100.0;
    }

    var stat = { x: date, y: cpuPercent };
    vm.cpuStats[0].values.push(stat);
    if (vm.cpuStats[0].values.length > 20) {
      vm.cpuStats[0].values.shift();
    }

  }

  function addMemoryUsage(date, usage) {
    var stat = { x: date, y: usage };
    vm.memoryStats[0].values.push(stat);
    if (vm.memoryStats[0].values.length > 20) {
      vm.memoryStats[0].values.shift();
    }
  }

  function addNetworkUsage(date, rxUsage, txUsage) {
    var rxStat = { x: date, y: rxUsage };
    vm.netStats[0].values.push(rxStat);
    if (vm.netStats[0].values.length > 20) {
      vm.netStats[0].values.shift();
    }

    var txStat = { x: date, y: txUsage };
    vm.netStats[1].values.push(txStat);
    if (vm.netStats[1].values.length > 20) {
      vm.netStats[1].values.shift();
    }
  }

  nv.addGraph(function() {
    graphCpuStats = nv.models.lineChart()
      .duration(0)
      .options({
        transitionDuration: 0,
        useInteractiveGuideline: true
      });

    graphCpuStats
      .x(function(d, i) { return d.x; });

    graphCpuStats.xAxis
      .tickFormat(function(d) { return d3.time.format('%H:%M:%S')(new Date(d)); })
      .axisLabel('');

    graphCpuStats.yAxis
      .axisLabel('%')
      .tickFormat(d3.format(',.2f'));

    graphCpuStats
      .forceY([0, 100])
      .showXAxis(true)
      .showYAxis(true);

    d3.select('#graphCpuStats')
      .datum(vm.cpuStats)
      .call(graphCpuStats);

    nv.utils.windowResize(graphCpuStats.update);

    return graphCpuStats;
  });

  nv.addGraph(function() {
    graphMemoryStats = nv.models.lineChart()
      .duration(0)
      .options({
        transitionDuration: 0
      });

    graphMemoryStats
      .x(function(d, i) { return d.x; });

    graphMemoryStats.xAxis
      .tickFormat(function(d) { return d3.time.format('%H:%M:%S')(new Date(d)); })
      .axisLabel('');

    graphMemoryStats.yAxis
      .axisLabel('MB')
      .tickFormat(d3.format(',.2f'));

    graphMemoryStats
      .forceY([0, 64])
      .showXAxis(true)
      .showYAxis(true);

    d3.select('#graphMemoryStats')
      .datum(vm.memoryStats)
      .call(graphMemoryStats);

    nv.utils.windowResize(graphMemoryStats.update);

    return graphMemoryStats;
  });

  nv.addGraph(function() {
    graphNetStats = nv.models.lineChart()
      .duration(0)
      .options({
        transitionDuration: 0
      });

    graphNetStats
      .x(function(d, i) { return d.x; });

    graphNetStats.xAxis
      .tickFormat(function(d) { return d3.time.format('%H:%M:%S')(new Date(d)); })
      .axisLabel('');

    graphNetStats.yAxis
      .axisLabel('MB')
      .tickFormat(d3.format(',.2f'));

    graphNetStats
    .forceY([0, 64])
      .showXAxis(true)
      .showYAxis(true);

    d3.select('#graphNetStats')
      .datum(vm.netStats)
      .call(graphNetStats);

    nv.utils.windowResize(graphNetStats.update);

    return graphNetStats;
  });

  var stream = oboe({
    url: '/containers/' + vm.id + '/stats',
    withCredentials: true,
    headers: {
      'Authorization': 'Bearer ' + store.get('sessionToken')
    }
  })
  .done(function(node) {
    // stats come every 1 second, only update according to the refresh interval
    var timestamp = Date.parse(node.read);
    var cpuCores = 1;
    if(node.cpu_stats.cpu_usage.percpu_usage) {
      cpuCores = node.cpu_stats.cpu_usage.percpu_usage.length;
    }
    addCpuUsage(timestamp, node.cpu_stats.system_cpu_usage, node.cpu_stats.cpu_usage.total_usage, cpuCores);
    // convert to MB
    addMemoryUsage(timestamp, node.memory_stats.usage / 1048576);
    // convert to MB
    var rxTotal = 0;
    var txTotal = 0;
    for(var iface in node.networks) {
      txTotal += node.networks[iface].tx_bytes / 1048576;
      rxTotal += node.networks[iface].rx_bytes / 1048576;
    }
    addNetworkUsage(timestamp, rxTotal, txTotal);
    refreshGraphs();
  })
  .fail(function(error) {
    MessageService.addErrorMessage('Error retrieving container stats', error.body);
  });

  function refreshGraphs() {
    graphCpuStats.update();
    graphMemoryStats.update();
    graphNetStats.update();
  }

  $scope.$on('$destroy', function() {
    stream.abort();
  });

  $('.ui.dropdown').dropdown();  // initiates Semantic UI dropdown
}

module.exports = ContainerStatsController;
