'use strict';

BannerController.$inject = ['$http', '$rootScope'];
function BannerController($http, $rootScope) {
  var vm = this;
  vm.messages = [];

  function getBanner() {
    $http({
      url: '/api/banner',
      noLoader: true
    }).success(function(data) {
      vm.messages = data;
    });
  }

  // Update the banner when loading the UI for the first time
  getBanner();

  // Update the banner every time a navigation event happens
  $rootScope.$on('$stateChangeSuccess', function() {
    getBanner();
  });
}

module.exports = BannerController;
