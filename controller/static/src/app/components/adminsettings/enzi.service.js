'use strict';

EnziService.$inject = ['$http', 'store'];
function EnziService($http, store) {
  return {
    getAuthBackend: function() {
      var promise = $http({
          method: 'GET',
          url: '/enzi/v0/config/auth'
        })
        .then(function(response) {
          return response.data.backend;
        });
      return promise;
    },
    setAuthBackend: function(backend) {
      return $http({
        method: 'PUT',
        url: '/enzi/v0/config/auth',
        data: {
          backend: backend
        }
      });
    },
    testLogin: function(username, password, config) {
      return $http({
        method: 'POST',
        url: '/enzi/v0/config/auth/ldap/tryLogin',
        data: {
          username: username,
          password: password,
          ldapSettings: config
        }
      });
    },
    getLdapConfig: function() {
      return $http({
        method: 'GET',
        url: '/enzi/v0/config/auth/ldap'
      });
    },
    updateLdapConfig: function(config) {
      return $http({
        method: 'PUT',
        url: '/enzi/v0/config/auth/ldap',
        data: config
      });
    },
    getJobs: function() {
      return $http({
        method: 'GET',
        url: '/enzi/v0/jobs?action=any&worker=any',
        noLoader: true
      });
    },
    triggerSync: function(config) {
      return $http({
        method: 'POST',
        url: '/enzi/v0/jobs',
        noLoader: true,
        data: {
          action: 'ldap-sync'
        }
      });
    },
    getJobLog: function(id) {
      return $http({
        method: 'GET',
        url: '/enzi/v0/jobs/' + id + '/logs',
				transformResponse: false
      });
    },
    getLegacyAuthSettings: function() {
      return $http.get('/api/config/auth2');
    },
    setLegacyAuthSettings: function(auth) {
      return $http({
        method: 'POST',
        url: '/api/config/auth2',
        data: auth
      });
    }
  };
}

module.exports = EnziService;
