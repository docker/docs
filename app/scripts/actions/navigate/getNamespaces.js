'use strict';

import _ from 'lodash';
import {
  Users
  } from 'hub-js-sdk';
var debug = require('debug')('navigate::getNamespaces');

export default function getNamespaces({actionContext, payload, done, maybeData}){
  var initialRepoName = payload.params.sourceRepoName;
  var initialNamespace = payload.location.query.namespace || maybeData.user.username;
  if (_.has(maybeData, 'token')) {
    Users.getNamespacesForUser(maybeData.token, function(err, res) {
      if (err) {
        return done();
      }
      Users.getUserSettings(maybeData.token, maybeData.user.username, function(settingsErr, settingsRes) {
        if (settingsRes.body) {
          actionContext.dispatch('RECEIVE_PRIVATE_REPOSTATS', settingsRes.body);
          actionContext.dispatch('RECEIVE_NAMESPACES', res.body);
          actionContext.dispatch('INITIALIZE_AUTOBUILD_FORM', {name: initialRepoName, namespace: initialNamespace});
          actionContext.dispatch('CLEAR_AUTOBUILD_FORM_ERRORS');
          if (payload.routes[payload.routes.length - 1].name === 'autobuildGithub') {
            actionContext.dispatch('SET_LINKED_REPO_TYPE', 'github');
          } else if (payload.routes[payload.routes.length - 1].name === 'autobuildBitbucket') {
            actionContext.dispatch('SET_LINKED_REPO_TYPE', 'bitbucket');
          }
          return done();
        } else if (settingsErr) {
          debug(settingsErr);
          return done();
        }
      });
    });
  } else {
    done();
  }
}
