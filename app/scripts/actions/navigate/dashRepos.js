'use strict';
var debug = require('debug')('navigate::dashRepos');
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
          actionContext.dispatch('CURRENT_USER_ORGS', res.body);
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

    //This is executed only for the currently logged in user, so put in starred, contributed only here
    var _doParallelCalls = function() {
      async.parallel([
        _getReposForUserOrOrg,
        _getOrgsForCurrentUser,
        _getUserSettings,
        function(callback) {
          Notifications.getActivityFeed(maybeData.token, function(err, res) {
            if (res) {
              actionContext.dispatch('RECEIVE_ACTIVITY_FEED', res.body);
            }
            callback();
          });
        }
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
        //get Org by name, if success set userOrOrg = `org` else set it to `user`
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
            actionContext.dispatch('CURRENT_USER_CONTEXT', { username: payload.params.user });
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
