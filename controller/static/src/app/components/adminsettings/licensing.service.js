'use strict';

LicensingService.$inject = ['$http'];
function LicensingService($http) {
  return {
    getLicense: function() {
      var promise = $http.get('/api/config/license')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    updateLicense: function(license) {
      return $http({
        method: 'POST',
        url: '/api/config/license',
        data: license
      });
    }
  };
}

module.exports = LicensingService;
