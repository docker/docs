'use strict';

FooterController.$inject = ['$http', 'store'];
function FooterController($http, store) {
  var vm = this;
  vm.ucpVersion = '';
  vm.ucpSHA = '';
  vm.apiVersion = '';
  $http.get('/version')
    .success(function(data) {
      if(data.Version) {
        vm.ucpVersion = data.Version.split('/')[1];
        store.set('ucpVersion', vm.ucpVersion);
      }
      vm.ucpSHA = data.GitCommit;
      vm.apiVersion = data.ApiVersion;
    });
}

module.exports = FooterController;
