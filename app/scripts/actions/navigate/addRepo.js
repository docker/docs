'use strict';

import has from 'lodash/object/has';
import {
  Users
} from 'hub-js-sdk';
import async from 'async';
var debug = require('debug')('navigate::ADD REPO');

export default function repo({actionContext, payload, done, maybeData}){

  // Defining functions here so we have access to payload
  const getUserSettings = (callback) => {
    // Users.getUserSettings === Orgs.getOrgSettings - equivalent calls
    let username = maybeData.user.username;
    if (payload.location.query.namespace) {
      username = payload.location.query.namespace;
    }
    Users.getUserSettings(maybeData.token, username, function(getErr, getRes) {
      // This is to get default visibility. If this fails, we shouldn't block
      if (getErr){
        callback(null, getErr);
      } else {
        callback(null, getRes.body);
      }
    });
  };

  const getNamespace = (callback) => {
    Users.getNamespacesForUser(maybeData.token, function(namespaceErr, namespaceRes) {
      if (namespaceErr) {
        // If we don't get back namespaces, we can't do anything, so no point in continuing with the other calls
        callback(namespaceErr);
      } else {
        callback(null, namespaceRes.body);
      }
    });
  };

  if (has(maybeData, 'token')) {
    async.parallel({
      getNamespace,
      getUserSettings
    },
    function(err, res) {
      actionContext.dispatch('CREATE_REPO_CLEAR_FORM');
      if (err) {
        done();
      } else {
        const is_private = has(res.getUserSettings, 'default_repo_visibility') ? res.getUserSettings.default_repo_visibility : true;
        actionContext.dispatch('CREATE_REPO_UPDATE_FIELD_WITH_VALUE', {fieldKey: 'is_private', fieldValue: is_private});
        actionContext.dispatch('RECEIVE_PRIVATE_REPOSTATS', res.getUserSettings);
        if (has(res.getNamespace, 'namespaces')) {
          actionContext.dispatch('CREATE_REPO_RECEIVE_NAMESPACES', {
            namespaces: res.getNamespace,
            selectedNamespace: maybeData.user.username
          });
        }
        // No namespaces is already handled in AddRepo.jsx - Dispatch is unecessary
        done();
      }
    });
  } else {
    done();
  }
}
