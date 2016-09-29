'use strict';

var debug = require('debug')('ACTION::saveOrgProfile');
import async from 'async';
import { Orgs } from 'hub-js-sdk';
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
export default function(actionContext, {jwt, orgname, organization}) {

  var _updateOrg = function(cb) {
    Orgs.updateOrg(jwt, orgname, organization, function(err, res) {
      if (err) {
        debug(err);
        actionContext.dispatch('UPDATE_ORG_ERROR', err);
        cb(err, {});
      } else {
        if (res.ok) {
          cb(null, res.body);
        }
      }
    });
  };

  //Get orgs for user
  var _getUpdatedOrg = function(cb) {
    Orgs.getOrg(jwt, orgname, function(err, res) {
      if (err) {
        debug(err);
        cb(err, {});
      } else {
        cb(null, res.body);
        actionContext.dispatch('UPDATE_ORG_SUCCESS');
      }
    });
  };

  async.series([
    _updateOrg,
    _getUpdatedOrg
  ], function (err, results) {
    if(err) {
      debug(err);
    } else {
      actionContext.dispatch('RECEIVE_ORGANIZATION', results[1]);
    }
  });
}
