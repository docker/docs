'use strict';

import { Builds } from 'hub-js-sdk';
const debug = require('debug')('hub:actions:linkGithub');

module.exports = function(actionContext, jwt) {
  Builds.getGithubClientID(jwt, function(err, res) {
    if (err) {
      debug('error', err);
      actionContext.dispatch('GITHUB_ID_ERROR', err);
    } else if (res.ok) {
      actionContext.dispatch('RECEIVE_GITHUB_ID', res.body);
    }
  });
};
