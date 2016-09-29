'use strict';

var debug = require('debug')('hub:actions:removeTeamMember');

import { Orgs } from 'hub-js-sdk';

var removeTeamMember = function(actionContext, {jwt, orgname, teamname, membername}) {
  Orgs.deleteMember(jwt, orgname, teamname, membername, function(delErr, delRes) {
    if (delErr) {
      debug('error', delErr);
      actionContext.dispatch('TEAM_MEMBER_ERROR', delErr);
    } else if (delRes.ok) {
      Orgs.getMembers(jwt, orgname, teamname, function(err, res) {
        if (err) {
          debug('getMembers error', err);
        } else {
          actionContext.dispatch('RECEIVE_DASHBOARD_TEAM_MEMBERS', res.body);
        }
      });
    }
  });
};

module.exports = removeTeamMember;
