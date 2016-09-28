'use strict';

var debug = require('debug')('hub:actions:associateGithubAccount');
import has from 'lodash/object/has';
import { Builds } from 'hub-js-sdk';

export default function(actionContext, {jwt, code}) {
  //TODO: immediately call the linked accounts status and update the authorized services page
  Builds.associateGithubAccount(jwt, code, function(err, res) {
    if (err) {
      debug('error', err);
      const { detail } = err.response.body;
      if(detail) {
        actionContext.dispatch('GITHUB_ASSOCIATE_ERROR', detail);
      }
    } else {
      if (res.body) {
        actionContext.dispatch('GITHUB_ASSOCIATE_SUCCESS', res.body);
      }
    }
  });
  actionContext.dispatch();
}
