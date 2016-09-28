'use strict';

LogoutController.$inject = ['AuthService', '$state', 'store'];
function LogoutController(AuthService, $state, store) {
  var vm = this;
  AuthService.logout();
  $state.go('login');
}

module.exports = LogoutController;
