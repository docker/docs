'use strict';

ResourcesController.$inject = ['AuthService'];
function ResourcesController(AuthService) {
  var vm = this;
  vm.account = AuthService.getCurrentUser();
}

module.exports = ResourcesController;
