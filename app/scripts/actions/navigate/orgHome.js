'use strict';
var debug = require('debug')('navigate orgHome');
import async from 'async';
import _ from 'lodash';
import {
  Repositories as Repos,
  Notifications,
  Orgs,
  Users
} from 'hub-js-sdk';

export default function home({actionContext, payload, done, maybeData}){
  if (_.has(maybeData, 'token')) {

    //Get repos for user or org
    var _getReposForUserOrOrg = function(cb, userType) {
      actionContext.dispatch('DASHBOARD_REPOS_STORE_ATTEMPTING_GET_REPOS');
      Repos.getReposForUser(maybeData.token, payload.params.user, function(err, res) {
        if (err) {
          actionContext.dispatch('ERROR_RECEIVING_REPOS');
          cb();
        } else {
          actionContext.dispatch('RECEIVE_REPOS', res.body);
          cb();
        }
      }, payload.location.query.page);
    };

    //Get contributed repos
    var _getReposByFilterTypeContrib = function(cb) {
      var username = maybeData.user.username;
      Repos.getContributedReposForUser(maybeData.token, username, function (err, res) {
        if (err) {
          cb();
        } else {
          var resultPayload = {type: 'contrib', results: res.body.results};
          actionContext.dispatch('RECEIVE_CONTRIB', resultPayload);
          cb();
        }
      });
    };

    //Get orgs for user
    var _getOrgsForCurrentUser = function(cb) {
      Users.getOrgsForUser(maybeData.token, function(err, res) {
        if (err) {
          debug(err);
          cb();
        } else {
          actionContext.dispatch('RECEIVE_DASHBOARD_NAMESPACES', {
            orgs: res.body,
            user: maybeData.user.username
          });
          cb();
        }
      });
    };

    //Get user settings for private repo stats
    var _getOrgPrivateRepoSettings = function(cb) {
      Orgs.getOrgSettings(maybeData.token, payload.params.user, function(err, res) {
        if (err) {
          debug(err);
          actionContext.dispatch('PRIVATE_REPOSTATS_NO_PERMISSIONS', err);
          actionContext.dispatch('TEAM_READ_ONLY', true);
          cb();
        } else {
          actionContext.dispatch('RECEIVE_PRIVATE_REPOSTATS', res.body);
          actionContext.dispatch('TEAM_READ_ONLY', false);
          cb();
        }
      });
    };

    //Get org's settings since we are in the organization's home
    var _getOrgSettings = function(cb) {
      Orgs.getOrg(maybeData.token, payload.params.user, function (err, res) {
        if (err) {
          cb();
        } else {
          actionContext.dispatch('RECEIVE_ORGANIZATION', res.body);
          cb();
        }
      });
    };

    var _getNamespaces = function(cb) {
      Users.getNamespacesForUser(maybeData.token, function(err, res) {
        if (err) {
          cb();
        } else {
          cb(null, res.body);
          actionContext.dispatch('CREATE_REPO_RECEIVE_NAMESPACES', {
            namespaces: res.body,
            selectedNamespace: maybeData.user.username
          });
        }
      });
    };

    //This is executed only for the currently logged in user, so put in starred, contributed only here
    var _doParallelCalls = function() {
      async.parallel([
        _getReposForUserOrOrg,
        _getOrgsForCurrentUser,
        _getReposByFilterTypeContrib,
        _getOrgSettings,
        _getOrgPrivateRepoSettings,
        _getNamespaces
      ], function(err, results) {
        if (err) {
          debug(err);
        }
        return done();
      });
    };

    //Get Organization by name
    var _getOrgByName = function(name, cb) {
      Orgs.getOrg(maybeData.token, name, function(err, res) {
        if (err) {
          debug(err);
          cb({});
        } else {
          var org = res.body;
          cb(org);
        }
      });
    };

    if (payload.params.user) {
      actionContext.dispatch('CURRENT_USER_CONTEXT', { username: payload.params.user });
      _doParallelCalls();
    } else {
      debug('mark 5');
      actionContext.dispatch('CURRENT_USER_CONTEXT', { username: payload.params.user });
      _doParallelCalls();
    }
  } else {
    debug('mark 6');
    done();
  }
}
