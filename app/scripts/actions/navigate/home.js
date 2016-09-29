'use strict';
var debug = require('debug')('navigate::home');
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
      var userOrOrgName = payload.params.user || maybeData.user.username;
      actionContext.dispatch('DASHBOARD_REPOS_STORE_ATTEMPTING_GET_REPOS');
      Repos.getReposForUser(maybeData.token, userOrOrgName, function(err, res) {
        if (err) {
          actionContext.dispatch('ERROR_RECEIVING_REPOS');
          cb();
        } else {
          actionContext.dispatch('RECEIVE_REPOS', res.body);
          cb();
        }
      }, payload.location.query.page);
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
    var _getUserSettings = function(cb) {
      Users.getUserSettings(maybeData.token, maybeData.user.username, function(err, res) {
        if (err) {
          debug(err);
          cb();
        } else {
          actionContext.dispatch('RECEIVE_PRIVATE_REPOSTATS', res.body);
          cb();
        }
      });
    };

    var _doParallelCalls = function() {
      async.parallel([
        _getReposForUserOrOrg,
        _getOrgsForCurrentUser,
        _getUserSettings
      ], function(err, results) {
        if (err) {
          debug(err);
        }
        return done();
      });
    };

    //Get Organization by name
    var _getOrgByName = function() {
      Orgs.getOrg(maybeData.token, payload.params.user, function(err, res) {
        if (err) {
          debug(err);
        } else {
          var org = res.body;
          if (!_.isEmpty(org)) {
            return org;
          }
        }
      });
      return {};
    };

    //Get user by name
    var _getUserByName = function() {
      Users.getUser(maybeData.token, payload.params.user, function(err, res) {
        if (err) {
          debug(err);
        } else {
          var user = res.body;
          if (!_.isEmpty(user)) {
            return user;
          }
        }
      });
      return {};
    };

    if (payload.params.user) {
      if (payload.params.user === maybeData.user.username) {
        actionContext.dispatch('CURRENT_USER_CONTEXT', { username: payload.params.user });
        _doParallelCalls();
      } else {
        async.series([
          function (callback) {
            callback(null, _getOrgByName(payload.params.user));
          },
          function (callback) {
            callback(null, _getUserByName(payload.params.user));
          }
        ], function (err, results) {
          if (err) {
            debug(err);
          } else {
            debug('Results of Series Call: ' + JSON.stringify(results));
            //User context is defined as `type`, `name` and `status flag`
            actionContext.dispatch('CURRENT_USER_CONTEXT', {username: payload.params.user });
            _getReposForUserOrOrg(done);
          }
        });
      }
    } else {
      actionContext.dispatch('CURRENT_USER_CONTEXT', { username: maybeData.user.username });
      _doParallelCalls();
    }
  } else {
    done();
  }
}
