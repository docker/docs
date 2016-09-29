'use strict';
var debug = require('debug')('navigate::dashStars');
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

    //This is always for current user
    //Get starred repos
    var _getReposByFilterTypeStarred = function(cb) {
      var username = maybeData.user.username;
      Repos.getStarredReposForUser(maybeData.token, username, function (err, res) {
        if (err) {
          cb();
        } else {
          actionContext.dispatch('RECEIVE_STARRED', res.body);
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

    //This is executed only for the currently logged in user, so put in starred, contributed only here
    var _doParallelCalls = function() {
      async.parallel([
        _getOrgsForCurrentUser,
        _getReposByFilterTypeStarred,
        _getUserSettings
      ], function(err, results) {
        if (err) {
          debug(err);
        }
        return done();
      });
    };

    actionContext.dispatch('CURRENT_USER_CONTEXT', { username: maybeData.user.username });
    _doParallelCalls();
  } else {
    done();
  }
}
