'use strict';

import _ from 'lodash';
import async from 'async';
import { PENDING_DELETE } from 'common/enums/RepoStatus';
import {
  Repositories as Repos
} from 'hub-js-sdk';

export default function repo({actionContext, payload, done, maybeData}){
  var token;
  if (_.has(maybeData, 'token')) {
    token = maybeData.token;
  } else {
    token = '';
  }
  var repoShortName = 'library/' + payload.params.splat;
  async.series([
    function(callback) {
      Repos.getRepo(token, repoShortName, function(err, res) {
        let status;
        if (res && res.body) {
          status = res.body.status;
        }

        if (err || status === PENDING_DELETE) {
          actionContext.dispatch('REPO_NOT_FOUND', err);
          callback(err);
        } else {
          actionContext.dispatch('RECEIVE_REPOSITORY', res.body);
          callback(null, res.body);
        }
      });
    }, function(callback) {
      Repos.getCommentsForRepo(token, repoShortName, function(err, res) {
        if (err) {
          callback(err);
        } else {
          actionContext.dispatch('RECEIVE_REPO_COMMENTS', res.body);
          callback(null, res.body);
        }
      });
    }
  ],
  function(error, results) {
    done();
  });
}
