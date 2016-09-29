'use strict';

import _ from 'lodash';
import {
  Orgs
  } from 'hub-js-sdk';
var debug = require('debug')('hub:actions:getTeamMembers');

export default function(actionContext, {jwt, orgname, teamname}) {
  Orgs.getMembers(jwt, orgname, teamname, function (err, res) {
    if (err) {
      debug('error', err);
    } else {
      actionContext.dispatch('RECEIVE_TEAM_MEMBERS', res.body);
    }
  });
}
