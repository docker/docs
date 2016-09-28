'use strict';

AccountsService.$inject = ['$http', '$q'];
function AccountsService($http, $q) {
  return {
    list: function() {
      var promise = $http
        .get('/api/accounts')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    listTeamMembers: function(teamId) {
      var promise = $http
        .get('/api/accounts?teamId=' + teamId)
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    getTeam: function(teamId) {
      var promise = $http
        .get('/api/teams/' + teamId)
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    deleteTeam: function(teamId) {
      var promise = $http
        .delete('/api/teams/' + teamId)
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    updateTeam: function(team) {
      var promise = $http
        .put('/api/teams/' + team.id, team)
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    addPermission: function(permission) {
      var promise = $http
        .post('/api/accesslists', permission)
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    getPermissions: function(teamId) {
      var promise = $http.get('/api/accesslists').
        then(function(response) {
          return response.data;
        });
      return promise;
    },
    getPermissionsForTeam: function(teamId) {
      var promise = $http
        .get('/api/accesslists')
        .then(function(response) {
          var accessLists = [];
          for(var i = 0; i < response.data.length; i++) {
            if(response.data[i].teamId === teamId) {
              accessLists.push(response.data[i]);
            }
          }
          return accessLists;
        });
      return promise;
    },
    getPermissionsForAccount: function(username) {
      // Get our list of permissions
      var promise = $http
        .get('/api/accesslists?username=' + username)
        .then(function(response) {

          // Create an API request per permission to grab the team name, since
          // there are many http requests, bundle them into an array of promises
          var promises = [];
          for(var i = 0; i < response.data.length; i++) {
            promises.push(
              $http.get('/api/teams/' + response.data[i].teamId)
            );
          }

          // Once all the promises are resolved, iterate through the responses
          return $q.all(promises)
              .then(function(teamResponses) {

                // Create a map of team ids -> team names
                var teamNames = {};
                for(i = 0; i < teamResponses.length; i++) {
                  teamNames[teamResponses[i].data.id] = teamResponses[i].data.name;
                }

                // Use the map of team ids -> team names to create an array of permissions
                var permissions = [];
                for(i = 0; i < response.data.length; i++) {
                  var permission = response.data[i];
                  permission.teamName = teamNames[permission.teamId];
                  permissions.push(permission);
                }

                return permissions;
          });
        });
      return promise;
    },
    deletePermission: function(permission) {
      var promise = $http
        .delete('/api/accesslists/' + permission.teamId + '/' + permission.id)
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    createUser: function(user) {
      var promise = $http
        .post('/api/accounts', user)
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    addUserToTeam: function(teamId, username) {
      var promise = $http
        .put('/api/teams/' + teamId + '/members/add/' + username)
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    removeUserFromTeam: function(teamId, username) {
      var promise = $http
        .delete('/api/teams/' + teamId + '/members/remove/' + username)
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    listTeams: function() {
      var promise = $http
        .get('/api/teams')
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    createTeam: function(team) {
      var promise = $http({
        method: 'POST',
        url: '/api/teams',
        data: team
      });
      return promise;
    },
    update: function(account) {
      var promise = $http({
        method: 'POST',
        url: '/api/accounts',
        headers: {
          'Content-Type': 'application/json'
        },
        data: account
      });
      return promise;
    },
    getAccount: function(username) {
      var promise = $http
        .get('/api/accounts/' + username)
        .then(function(response) {
          return response.data;
        });
      return promise;
    },
    removeAccount: function(account) {
      var promise = $http
        .delete('/api/accounts/' + account.username)
        .then(function(response) {
          return response.data;
        });
      return promise;
    }
  };
}

module.exports = AccountsService;
