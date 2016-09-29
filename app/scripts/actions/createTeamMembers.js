'use strict';

var debug = require('debug')('hub:actions:createTeamMembers');
import async from 'async';
import { Orgs } from 'hub-js-sdk';

export default function(actionContext, {jwt, orgname, teamname, members}) {

  var _createMember = function(cb) {
    var _makeMemberFuncArray = function(member) {
      return function() {
        // member.username can be either a valid email or a username
        Orgs.createMember(jwt, orgname, teamname, member.username, function(err, res) {
          if (err) {
            debug('createMember error', err);
            if (res.badRequest) {
              actionContext.dispatch('TEAM_MEMBER_BAD_REQUEST', err);
            } else if (res.unauthorized) {
              actionContext.dispatch('TEAM_MEMBER_UNAUTHORIZED', err);
            } else {
                actionContext.dispatch('TEAM_MEMBER_ERROR', err);
            }
          } else {
            if (res.ok) {
              actionContext.dispatch('ORG_TEAM_CLEAR_ERROR_STATES');
              cb(null, res.body);
            }
          }
        });
      };
    };
    var funcArray = [];
    for(var i = 0; i < members.length; ++i) {
      funcArray.push(_makeMemberFuncArray(members[i]));
    }
    async.parallel(funcArray, function(err, results) {
      if (err) {
        cb(err);
      } else {
        cb(null, results);
      }
    });
  };

  //Get members for team
  var _getMembersForTeam = function(cb) {
    Orgs.getMembers(jwt, orgname, teamname, function(err, res) {
      if (err) {
        debug('getMembers error', err);
        cb(err, {});
      } else {
        actionContext.dispatch('RECEIVE_DASHBOARD_TEAM_MEMBERS', res.body);
        cb(null, res.body.results);
      }
    });
  };

  async.series([
    _createMember,
    _getMembersForTeam
  ], function (err, results) {
    if(err) {
      debug('final callback error', err);
    }
  });
}
