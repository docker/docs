'use strict';

var $ = require('jquery');

LoginController.$inject = ['AuthService', 'LicensingService', '$state', '$stateParams', '$location', '$http', '$rootScope'];
function LoginController(AuthService, LicensingService, $state, $stateParams, $location, $http, $rootScope) {
  var vm = this;
  vm.error = '';
  vm.username = '';
  vm.password = '';
  vm.login = login;

  var returnLocation = $stateParams.return;

  function isValid() {
    return $('#login .ui.form').form('validate form');
  }

  function isUserAnAdmin() {
    return AuthService.getMyAccount()
      .then(function(data) {
        return data.data.admin;
      });
  }

  function isUcpLicensed() {
    return $http.get('/info')
      .then(function(data) {
        for(var i = 0; i < data.data.Labels.length; i++) {
          if(data.data.Labels[i] === 'com.docker.ucp.license_key=unlicensed') {
            return false;
          }
        }
        return true;
      }, function(error) {
        return true;
      });
  }

  function redirectSuccessfulLogin() {
    if($rootScope.returnState) {
      // If the return parameter is set, go to requested page
      $state.go($rootScope.returnState, $rootScope.returnStateParams);
      $rootScope.returnState = '';
      $rootScope.returnParams = {};
    } else {
      // Otherwise go to dashboard
      $state.go('dashboard.main');
    }

  }

  function login() {
    if (!isValid()) {
      return;
    }

    vm.error = '';

    AuthService.login({
      username: vm.username,
      password: vm.password
    }).then(function(data) {

      isUserAnAdmin()
        .then(function(admin) {
          if(admin === true) {
            // If the user is an admin, and the UCP install is unlicensed, go to the dashboard
            isUcpLicensed()
              .then(function(licensed) {
                if(licensed === true) {
                  redirectSuccessfulLogin();
                } else {
                  $state.go('licenseReminder');
                }
              });
          } else {
            redirectSuccessfulLogin();
          }
        });

    }, function(error) {
      var resp = error.data.trim();
      if (resp === 'unauthorized') {
        vm.error = 'Incorrect authentication credentials.  This may also be due to an invalid client certificate.';
      } else {
        vm.error = error.data;
      }
    });
  }
}

module.exports = LoginController;
