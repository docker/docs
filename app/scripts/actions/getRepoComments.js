'use strict';

import _ from 'lodash';
import {
  Repositories as Repos
  } from 'hub-js-sdk';

export default function repoComments(actionContext, payload) {
  var token = payload.JWT;

  var namespace = payload.namespace;
  if(payload.namespace === '_') {
    namespace = 'library';
  }
  var repoShortName = namespace + '/' + payload.repoName;
  Repos.getCommentsForRepo(token, repoShortName, function(err, res) {
    if (err) {
      actionContext.dispatch('ERROR_RECEIVING_REPO_COMMENTS');
    } else {
      actionContext.dispatch('RECEIVE_REPO_COMMENTS', res.body);
    }
  }, payload.pageNumber);
}
