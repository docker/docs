'use strict';

var debug = require('debug')('hub:actions:deleteRepoComment');

import { Repositories } from 'hub-js-sdk';

var deleteRepoComment = function(actionContext, {jwt, repoShortName, commentid}) {
  Repositories.deleteRepoComment(jwt, repoShortName, commentid, function(delErr, delRes) {
    if (delErr) {
      debug('deleteRepoComment error', delErr);
    } else if (delRes.ok) {
      Repositories.getCommentsForRepo(jwt, repoShortName, function(getErr, getRes) {
        if (getErr) {
          debug('getCommentsForRepo error', getErr);
        } else {
          actionContext.dispatch('RECEIVE_REPO_COMMENTS', getRes.body);
        }
      });
    }
  });

};

module.exports = deleteRepoComment;
