'use strict';

var STORE_AUTH_USERNAME = 'username';
var STORE_AUTH_ADMIN = 'userIsAdmin';
var STORE_AUTH_ROLE = 'userRole';
var STORE_AUTH_SESSION_TOKEN = 'sessionToken';

AuthService.$inject = ['$http', '$state', '$rootScope', '$timeout', '$q', 'store'];
function AuthService($http, $state, $rootScope, $timeout, $q, store) {
  function logout() {
    return $http.post('/auth/logout', true)
      .then(function(data) {
        store.remove(STORE_AUTH_SESSION_TOKEN);
      });
  }

  function updateRootScopeAccountInfo() {
    $http({
      method: 'GET',
      url: '/api/account'
    }).then(function(data) {
      store.set(STORE_AUTH_USERNAME, data.data.username);
      store.set(STORE_AUTH_ADMIN, data.data.admin);
      store.set(STORE_AUTH_ROLE, data.data.role);

      // TODO: Remove contextual data from root scope, this is used directly in
      // the HTML templates at the moment
      $rootScope.username = data.data.username;
      $rootScope.userIsAdmin = data.data.admin;
      $rootScope.userRole = data.data.role;

      $timeout(function() {
        $rootScope.$apply();
      });
    });
  }
  return {
    getCurrentUser: function() {
      return {
        username: store.get(STORE_AUTH_USERNAME),
        admin: store.get(STORE_AUTH_ADMIN),
        role: store.get(STORE_AUTH_ROLE),
        sessionToken: store.get(STORE_AUTH_SESSION_TOKEN)
      };
    },
    login: function(credentials) {
      var loggedOut = $q.defer();

      // If a sessionToken already exists in localStorage then attempt log out
      if(store.get(STORE_AUTH_SESSION_TOKEN)) {
        logout().finally(function() {
          loggedOut.resolve();
        });
      } else {
        // If not, resolve the promise and continue
        loggedOut.resolve();
      }

      return loggedOut.promise.then(function() {
        return $http({
            method: 'POST',
            url: '/auth/login',
            data: credentials,
            ignore401: true
          })
          .then(function(data) {
            store.set(STORE_AUTH_SESSION_TOKEN, data.data.auth_token);
            updateRootScopeAccountInfo();
          });
      });
    },
    logout: logout,
    updateRootScopeAccountInfo: updateRootScopeAccountInfo,
    getMyAccount: function() {
      var promise = $http({
        method: 'GET',
        url: '/api/account'
      });
      return promise;
    }
  };
}

module.exports = AuthService;
