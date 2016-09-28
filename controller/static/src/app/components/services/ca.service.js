'use strict';

CAService.$inject = ['$http'];
function CAService($http) {
  return {
    getCACert: function() {
      return $http({
          method: 'GET',
          url: '/ca'
        })
        .then(function(data) {
          return data.data;
        }, function(error) {
          return '';
        });
    }
  };
}

module.exports = CAService;
