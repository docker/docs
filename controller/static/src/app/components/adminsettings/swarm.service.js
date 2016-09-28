'use strict';

SwarmService.$inject = ['$http'];
function SwarmService($http) {
  return {
    getSwarmConfig: function() {
      return $http({
        method: 'GET',
        url: '/swarm'
      });
    },
    updateSwarm: function(config, version) {
      return $http({
        method: 'POST',
        url: '/swarm/update?version=' + version,
        data: config
      });
    }
  };
}

module.exports = SwarmService;
