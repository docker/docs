'use strict';

import { Users } from 'hub-js-sdk';
var debug = require('debug')('hub:actions:toggleVisibility');

module.exports = function(actionContext, {JWT, username, visibility}) {
  actionContext.dispatch('START TOGGLE REPO VISIBILITY');
  Users.updateDefaultVisibility(JWT, { username, visibility }, function(err, res){
    if (err) {
      actionContext.dispatch('UPDATE_ORG_ERROR', err);
    } else if (res.ok) {
      Users.getUserSettings(JWT, username, function(getErr, getRes) {
        if (getErr) {
          debug(getErr);
        } else if (getRes.ok) {
          actionContext.dispatch('RECEIVE_PRIVATE_REPOSTATS', getRes.body);
          let is_private = (getRes.body.default_repo_visibility === 'private');
          actionContext.dispatch('CREATE_REPO_UPDATE_FIELD_WITH_VALUE', {fieldKey: 'is_private', fieldValue: is_private});
        }
      });
    }
  });
};
