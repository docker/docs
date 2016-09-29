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
    const currentOrg = payload.params.user;

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
      Orgs.getOrgSettings(maybeData.token, currentOrg, function(err, res) {
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

    //Get Teams for an org
    var _getTeamsForOrg = function(cb) {
      Orgs.getTeams(maybeData.token, currentOrg, function(err, res) {
        if (err) {
          debug(err);
          cb();
        } else {
          actionContext.dispatch('RECEIVE_DASHBOARD_ORG_TEAMS', res.body);
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

    const currentTeam = payload.location.query.team;
    var _getTeam = function(cb) {
      if (currentTeam) {
        Orgs.getTeam(maybeData.token, currentOrg, currentTeam, function(err, res) {
          if (err) {
            actionContext.dispatch('ORG_DASHBOARD_MEMBERS_ERROR', err);
            cb(null, err);
          } else {
            actionContext.dispatch('RECEIVE_DASHBOARD_ORG_TEAM', res.body);
            cb(null, res.body);
          }
        });
      } else {
        cb(null, null);
      }
    };

    var _getMembers = function(cb) {
      if (currentTeam) {
        Orgs.getMembers(maybeData.token, currentOrg, currentTeam, function(err, res) {
          if (err) {
            actionContext.dispatch('ORG_DASHBOARD_MEMBERS_ERROR', err);
            cb(null, err);
          } else {
            actionContext.dispatch('RECEIVE_DASHBOARD_TEAM_MEMBERS', res.body);
            cb(null, res.body);
          }
        });
      } else {
        cb(null, null);
      }
    };

    //This is executed only for the currently logged in user, so put in starred, contributed only here
    var _doParallelCalls = function() {
      async.parallel([
        _getOrgsForCurrentUser,
        _getOrgPrivateRepoSettings,
        _getTeamsForOrg,
        _getNamespaces,
        _getTeam,
        _getMembers
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

    actionContext.dispatch('CURRENT_USER_CONTEXT', {
      username: payload.params.user
    });
    _doParallelCalls();

  } else {
    done();
  }
}
