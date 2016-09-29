'use strict';
import { Orgs } from 'hub-js-sdk';
import async from 'async';
const debug = require('debug')('hub:actions:createOrgTeam');

module.exports = function(actionContext, {jwt, orgName, team}) {

  var _createTeam = function(cb) {
    Orgs.createTeam(jwt, orgName, team, function(err, res) {
      if (err) {
        debug('createTeam error', err);
        cb(err);
        var errRes = err.response;
        if (errRes.badRequest) {
          actionContext.dispatch('TEAM_BAD_REQUEST', err);
        } else if (errRes.unauthorized) {
          actionContext.dispatch('TEAM_UNAUTHORIZED', err);
        } else {
          actionContext.dispatch('TEAM_ERROR', err);
        }
      } else {
        if (res.body) {
          //send team name to get members for the team
          actionContext.dispatch('CREATE_ORG_TEAM', res.body);
          actionContext.dispatch('RECEIVE_DASHBOARD_ORG_TEAM', res.body);
          cb(null, res.body.name);
        }
      }
    });
  };

  var _getMembers = function(teamName, cb) {
    Orgs.getMembers(jwt, orgName, teamName, function(err, res) {
      if (err) {
        debug('getMembers error', err);
        cb(err);
      } else {
        actionContext.dispatch('RECEIVE_DASHBOARD_TEAM_MEMBERS', res.body);
        cb(null, 'done');
      }
    });
  };

  var _createTeamWrapper = function(cb) {
    async.waterfall([
      _createTeam,
      _getMembers
      ], function(err, res) {
        cb(null, res);
      }
    );
  };

  var _getTeamsForOrg = function(cb) {
    Orgs.getTeams(jwt, orgName, function(err, res) {
      if (err) {
        debug('getTeams error', err);
        cb(err);
      } else {
        if (res.body) {
          cb(null, res.body);
        }
      }
    });
  };

  async.series([
    _createTeamWrapper,
    _getTeamsForOrg
  ], function(err, results) {
    if (results[0] && results[1]) {
      actionContext.dispatch('RECEIVE_DASHBOARD_ORG_TEAMS', results[1]);
      actionContext.history.push(`/u/${orgName}/dashboard/teams/?team=${team.name}`);
    }
  });
};
