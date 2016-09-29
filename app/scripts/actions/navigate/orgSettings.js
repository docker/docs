'use strict';

import _ from 'lodash';
import {
  Users
  } from 'hub-js-sdk';
import async from 'async';
var debug = require('debug')('navigate::orgSettings');

export default function orgs({actionContext, payload, done, maybeData}){
  debug('Organization Settings Navigate Token -> ' + maybeData.token);
  debug('ORG SETTINGS PAYLOAD', payload);
  let orgName = payload.params.user;
  if (_.has(maybeData, 'token')) {
    async.parallel({
      getOrgs: function(callback) {
        Users.getOrgsForUser(maybeData.token, function(getErr, getRes) {
          if (getErr) {
            debug(getErr);
            callback();
          } else {
            actionContext.dispatch('RECEIVE_DASHBOARD_NAMESPACES', {
              orgs: getRes.body,
              user: maybeData.user.username
            });
            actionContext.dispatch('CURRENT_USER_ORGS', getRes.body);
            actionContext.dispatch('SELECT_ORGANIZATION', orgName);
            actionContext.dispatch('CURRENT_USER_CONTEXT', { username: orgName });
            callback(null, getRes.body);
          }
        });
      },
      getUserSettings: function(callback) {
        let username = payload.params.user;
        Users.getUserSettings(maybeData.token, username, function(getErr, getRes) {
          if (getErr){
            debug(getErr);
            callback();
          } else {
            let is_private = (getRes.default_repo_visibility);
            actionContext.dispatch('CREATE_REPO_UPDATE_FIELD_WITH_VALUE', {fieldKey: 'is_private', fieldValue: is_private});
            actionContext.dispatch('RECEIVE_PRIVATE_REPOSTATS', getRes.body);
            callback(null, getRes.body);
          }
        });
      },
      getNamespaces: function(callback) {
        Users.getNamespacesForUser(maybeData.token, function(err, res) {
          if (err) {
            callback();
          } else {
            callback(null, res.body);
            actionContext.dispatch('CREATE_REPO_RECEIVE_NAMESPACES', {
              namespaces: res.body,
              selectedNamespace: maybeData.user.username
            });
          }
        });
      }
    }, function(err, res){
      done();
    });
  } else {
    done();
  }
}
