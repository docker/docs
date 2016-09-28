/* @flow */
'use strict';

import async from 'async';
import { Repositories as Repos } from 'hub-js-sdk';
import { PENDING_DELETE } from 'common/enums/RepoStatus';
var debug = require('debug')('hub:actions:attemptChangeShortDescription');

export default function(actionContext,
                        {jwt, repoShortName, shortDescription},
                        done) {
  actionContext.dispatch('SHORT_DESCRIPTION_ATTEMPT_START');

  var name = repoShortName.split('/')[1];
  var namespace = repoShortName.split('/')[0];
  var _updateShortDescription = function(cb) {
    Repos.patchRepo(jwt, repoShortName, {description: shortDescription},
      function(err, res) {
        if (err) {
          if(res && res.badRequest) {
            debug('error', err);
            actionContext.dispatch('SHORT_BAD_REQUEST', res.body);
            cb(err);
          } else {
            actionContext.dispatch('DETAILS_ERROR');
            cb(err);
          }
        } else {
          actionContext.dispatch('SHORT_DESCRIPTION_SUCCESS');
          cb(null, res.body);
        }
      }
    );
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
    _updateShortDescription,
    _getRepoDetails
  ], function (err, results) {
    if(err) {
      debug('error', err);
    }
  });
}
