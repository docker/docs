'use strict';

import { Builds } from 'hub-js-sdk';
const debug = require('debug')('hub:actions:linkBitbucket');

module.exports = function(actionContext, jwt) {
  Builds.getBitbucketAuthUrl(jwt, function(err, res) {
    if (err) {
      debug('error', err);
      actionContext.dispatch('BITBUCKET_AUTH_URL_ERROR', err);
    } else if (res.ok) {
      actionContext.dispatch('RECEIVE_BITBUCKET_AUTH_URL', res.body);
    }
  });
};
