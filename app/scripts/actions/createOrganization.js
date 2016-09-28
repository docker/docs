'use strict';

var debug = require('debug')('hub:actions:createOrganization');
import async from 'async';
import { Orgs, Users } from 'hub-js-sdk';
//Organization Object
/*
 {
 id (string),
 orgname (regex),
 full_name (string),
 location (string): Your Location on the world,
 company (string): Your organization's name,
 profile_url (url): Your place on the web,
 gravatar_email (email): This address will define which picture of you is shown,
 is_active (boolean): Designates whether user is active. Unselect this instead of deleting accounts.,
 date_joined (datetime),
 gravatar_url (string)
 }
 */
export default function(actionContext, {jwt, organization}) {

  var _createOrg = function(cb) {
    Orgs.createOrg(jwt, organization, function(err, res) {
      if (err) {
        if(res && res.badRequest) {
          debug('createOrg error', err);
          actionContext.dispatch('ADD_ORG_BAD_REQUEST', res.body);
          cb(err, res);
        } else {
          actionContext.dispatch('ADD_ORG_FACEPALM');
          cb(err, res);
        }
      } else {
        cb(null, res.body);
      }
    });
  };

  //Get orgs for user
  var _getOrgsForUser = function(cb) {
    Users.getOrgsForUser(jwt, function(err, res) {
      if (err) {
        debug('getOrgsForUser error', err);
        cb(err, {});
      } else {
        actionContext.dispatch('CURRENT_USER_ORGS', res.body.results);
        cb(null, res.body.results);
      }
    });
    return {};
  };

  async.series([
      _createOrg,
      _getOrgsForUser
  ], function (err, results) {
       if(err) {
         debug('final callback error', err);
       } else {
         actionContext.dispatch('CREATED_ORGANIZATION', {newOrg: results[0], userOrgs: results[1]});
         actionContext.history.push(`/u/${organization.orgname}/dashboard/teams/`);
       }
  });
}
