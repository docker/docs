'use strict';

var debug = require('debug')('hub:actions:removeTeam');

import { Orgs } from 'hub-js-sdk';

var removeTeam = function(actionContext, {jwt, orgname, teamname}) {
  Orgs.deleteTeam(jwt, orgname, teamname, function(delErr, delRes) {
    if (delErr) {
      debug('error', delErr);
      actionContext.dispatch('TEAM_ERROR', delErr);
    } else if (delRes.ok) {
      actionContext.dispatch('DELETE_DASHBOARD_TEAM_SUCCESS');
      actionContext.history.push(`/u/${orgname}/dashboard/teams/`);
      Orgs.getTeams(jwt, orgname, function(err, res) {
        if (err) {
          debug('getTeams error', err);
        } else {
          actionContext.dispatch('RECEIVE_DASHBOARD_ORG_TEAMS', res.body);
        }
      });
    }
  });
};

module.exports = removeTeam;
