'use strict';

import async from 'async';
var debug = require('debug')('hub:actions:updateRepoSettingsVisibilityField');
import { Repositories as Repos, Users } from 'hub-js-sdk';

function updateRepoVisibility(actionContext, { jwt, isPrivate, repoShortName }) {
  return (cb) => {
    Repos.updateRepoVisibility(jwt, { isPrivate, repoShortName }, (err, res) => {
      if (err) {
        const { badRequest, body } = err.response;
        if(badRequest) {
          debug(err);
          actionContext.dispatch('VISIBILITY_BAD_REQUEST', body);
          cb(err);
        } else if (res && res.body.error) {
          actionContext.dispatch('VISIBILITY_ERROR', res.body);
          cb(err);
        } else {
          actionContext.dispatch('VISIBILITY_ERROR');
          cb(err);
        }
      } else {
        actionContext.dispatch('TOGGLE_VISIBILITY_SUCCESS', isPrivate);
        cb(null, res.body);
      }
    });
  };
}

function getPrivateRepoStats(actionContext, { jwt, repoShortName }) {
  return (cb) => {
    const namespace = repoShortName.split('/')[0];
    Users.getUserSettings(jwt, namespace, function (err, res) {
      if (err) {
        cb();
      } else {
        actionContext.dispatch('RECEIVE_PRIVATE_REPOSTATS', res.body);
        cb();
      }
    });
  };
}


export default function(actionContext, { jwt, isPrivate, repoShortName }, done) {
  actionContext.dispatch('TOGGLE_VISIBILITY_ATTEMPT_START');
  async.series([
    updateRepoVisibility(actionContext, { jwt, isPrivate, repoShortName }),
    getPrivateRepoStats(actionContext, {jwt, repoShortName })
  ], function(err, results) {
    if (err) {
      debug(err);
    }
  });
}
