'use strict';

var debug = require('debug')('hub:actions:saveTeamProfile');
import async from 'async';
import { Orgs } from 'hub-js-sdk';
//Team Object
/*
 {
 id (string),
 teamname (regex),
 description (string)
 }
 */
export default function(actionContext, {jwt, orgname, teamname, team}) {

  var _updateTeam = function(cb) {
    Orgs.updateTeam(jwt, { orgname, teamname, team }, function(err, res) {
      if (err) {
        debug(err);
        actionContext.dispatch('UPDATE_TEAM_ERROR', err);
        cb(err, {});
      } else {
        if (res.ok) {
          actionContext.dispatch('UPDATE_TEAM_SUCCESS');
          cb(null, res.body);
        }
      }
    });
  };

  //Get Team for org
  var _getUpdatedTeam = function(cb) {
    Orgs.getTeam(jwt, orgname, team.name, function(err, res) {
      if (err) {
        debug(err);
        cb(err, {});
      } else {
        cb(null, res.body);
      }
    });
  };

  async.series([
    _updateTeam,
    _getUpdatedTeam
  ], function (err, results) {
    if(err) {
      debug(err);
    } else {
      actionContext.dispatch('RECEIVE_ORG_TEAM', results[1]);
    }
  });
}
