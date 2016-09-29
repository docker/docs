'use strict';

import { Repositories as Repos } from 'hub-js-sdk';
const debug = require('debug')('hub:actions:deleteRepo');

export default function(actionContext, {jwt, repoShortName}) {
  actionContext.dispatch('DELETE_REPO_ATTEMPT_START');
  Repos.deleteRepository(jwt, repoShortName, (err, res) => {
    if (err) {
      const { badRequest, body } = err.response;
      if(badRequest) {
        actionContext.dispatch('DELETE_REPO_BAD_REQUEST', body);
      } else {
        actionContext.dispatch('DELETE_REPO_ERROR', err);
      }
    } else {
      if (res && res.ok) {
        actionContext.history.push('/');
      }
    }
  });
}
