'use strict';

var angular = require('angular');

DashboardService.$inject = ['$http'];
function DashboardService($http) {
  return {
    info: function() {
      var promise = $http
      .get('/info')
      .then(function(response) {
        return response.data;
      });
      return promise;
    },
    nodeInfo: function() {
      return $http
        .get('/nodes')
        .then(function(response) {
          return response.data;
        });
    },
    containerInfo: function() {
      return $http
        .get('/containers/json?all=1')
        .then(function(response) {
          return response.data;
        });
    },
    parseSystemStatus: function(systemStatus) {
      if(!systemStatus) {
        return {};
      }
      var res = {};
      res.nodeCount = 0;
      res.reservedCpuTotal = 0;
      res.reservedCpuUsed = 0;
      res.reservedCpu = 0;
      res.reservedMemTotal = 0;
      res.reservedMemUsed = 0.0;
      res.strategy = 'Unknown';
      res.strategyIconClass = 'help icon';
      res.controllers = [];

      for(var i = 0; i < systemStatus.length; i++) {
        var item = systemStatus[i];
        if(item[0].localeCompare('Strategy') === 0) {
          res.strategy = item[1];
          if(res.strategy.localeCompare('spread') === 0) {
            res.strategyIconClass = 'maximize icon';
          } else if(res.strategy.localeCompare('binpack') === 0) {
            res.strategyIconClass = 'compress icon';
          } else if(res.strategy.localeCompare('random') === 0) {
            res.strategyIconClass = 'random icon';
          }
        } else if(item[0].localeCompare('Filters') === 0) {
          res.filters = item[1];
        } else if(item[0].localeCompare('Nodes') === 0) {
          res.nodeCount = parseInt(item[1]);
        } else if(item[0].indexOf('Reserved CPUs') !== -1) {
          var parts = item[1].split(' ');
          res.reservedCpuUsed += parseFloat(parts[0]);
          res.reservedCpuTotal += parseFloat(parts[2]);
        } else if(item[0].indexOf('Reserved Memory') !== -1) {
          var memparts = item[1].split(' ');
          res.reservedMemUsed += normalizeMemoryToGiB(parseFloat(memparts[0]), memparts[1]);
          res.reservedMemTotal += normalizeMemoryToGiB(parseFloat(memparts[3]), memparts[4]);
        } else if(item[0].indexOf('Cluster Managers') !== -1) {
          var managerCount = parseInt(item[1]);
          i++;
          for(var j = 0; j < managerCount; j++) {
            // FIXME: This parsing is super fragile, will not cater for future modifications at all
            var controller = {
              status: systemStatus[i][1],
              url: systemStatus[i + 1][1],
              manager: systemStatus[i + 2][1],
              kv: systemStatus[i + 3][1]
            };
            res.controllers.push(controller);
            i += 4;
          }
        }
      }

      // Calculate progression percentages
      res.reservedMemCalc = Math.ceil(res.reservedMemUsed / res.reservedMemTotal * 100);
      res.reservedCpuCalc = Math.ceil(res.reservedCpuUsed / res.reservedCpuTotal * 100);

      // Defensively set 100 to be the maximum
      res.reservedMem = Math.min(res.reservedMemCalc, 100);
      res.reservedCpu = Math.min(res.reservedCpuCalc, 100);
      return res;
    }
  };
}

function normalizeMemoryToGiB(value, measurement) {
  if(measurement === 'MiB') {
    return value / 1000.0;
  } else if (measurement === 'KiB') {
    return value / 1000000.0;
  } else if (measurement === 'B') {
    return value / 1000000000.0;
  } else if (measurement === 'GiB') {
    return value;
  } else if (measurement === 'TiB') {
    return value * 1000.0;
  } else {
    return 0;
  }
}



module.exports = DashboardService;
