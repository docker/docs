/* @flow */
'use strict';

import { Repositories as Repos } from 'hub-js-sdk';
import async from 'async';
import { PENDING_DELETE } from 'common/enums/RepoStatus';
var debug = require('debug')('hub:actions:attemptChangeLongDescription');

export default function(actionContext,
                        {jwt, repoShortName, longDescription},
                        done) {
  actionContext.dispatch('LONG_DESCRIPTION_ATTEMPT_START');

  var _updateLongDescription = function(cb) {
    Repos.patchRepo(jwt, repoShortName, {
      full_description: longDescription
    }, function(err, res) {
      if (err) {
        if(res && res.badRequest) {
          debug('error', err);
          actionContext.dispatch('LONG_BAD_REQUEST', res.body);
          cb(err);
        } else {
          actionContext.dispatch('DETAILS_ERROR');
          cb(err);
        }
      } else {
        actionContext.dispatch('LONG_DESCRIPTION_SUCCESS');
        cb(null, res.body);
      }
    });
  };

  var _getRepoDetails = function(cb) {
    Repos.getRepo(jwt, repoShortName, function(err, res) {
      const { status, detail } = res.body;

      if (err || status === PENDING_DELETE) {
        actionContext.dispatch('REPO_NOT_FOUND', err);
        cb(null, detail);
      } else {
        actionContext.dispatch('RECEIVE_REPOSITORY', res.body);
        cb(null, res.body);
      }
    });
  };

  async.series([
    _updateLongDescription,
    _getRepoDetails
  ], function (err, results) {
    if(err) {
      debug('error', err);
    }
  });
}
