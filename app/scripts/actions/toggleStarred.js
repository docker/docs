'use strict';

import {
  Repositories as Repos
} from 'hub-js-sdk';
var debug = require('debug')('hub:actions:toggleStarred');

export default function(actionContext, {jwt, repoShortName, status}) {
  if(status) {
    Repos.starRepo(jwt, repoShortName,
      function(err, res) {
        if (err) {
          debug(err);
        } else if (res.ok) {
          actionContext.dispatch('TOGGLE_STARRED_STATE', status);
        }
      }
    );
  } else {
    Repos.unstarRepo(jwt, repoShortName,
      function(err, res) {
        if (err) {
          debug(err);
        } else if (res.ok) {
          actionContext.dispatch('TOGGLE_STARRED_STATE', status);
        }
      }
    );
  }
}
