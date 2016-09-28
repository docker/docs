'use strict';

var debug = require('debug')('hub:actions:addRepoComment');

import { Repositories } from 'hub-js-sdk';

var addRepoComment = function(actionContext, {jwt, repoShortName, comment}) {
  Repositories.addCommentToRepo(jwt, repoShortName, comment, function(addErr, addRes) {
    if (addErr) {
      debug(addErr);
    } else if (addRes.body && addRes.ok) {
      Repositories.getCommentsForRepo(jwt, repoShortName, function(getErr, getRes) {
        if (getErr) {
          debug(getErr);
        } else {
          actionContext.dispatch('RECEIVE_REPO_COMMENTS', getRes.body);
        }
      }, 1);
    }
  });

};

module.exports = addRepoComment;
