'use strict';
const debug = require('debug')('navigate::repo');
import { parallel, waterfall } from 'async';
import { PENDING_DELETE } from 'common/enums/RepoStatus';
import _ from 'lodash';
import {
  Repositories as Repos,
  Users
} from 'hub-js-sdk';

function getRepo({maybeToken, actionContext, user, splat}) {
  return function(callback){
    Repos.getRepo(maybeToken, `${user}/${splat}`, function(err, res) {
      let status;
      if (res && res.body) {
        status = res.body.status;
      }

      if (err || status === PENDING_DELETE) {
        actionContext.dispatch('REPO_NOT_FOUND', res.body);
        return callback(null, null);
      } else {
        actionContext.dispatch('RECEIVE_REPOSITORY', res.body);
        return callback(null, res.body);
      }
     });
  };
}

function handleGetPrivateRepoStats({maybeToken, user, actionContext}) {
  return function(repoDetails, callback) {
    if (repoDetails) {
      Users.getUserSettings(maybeToken, user, function (err, res) {
        if (err) {
          callback();
        } else {
          actionContext.dispatch('RECEIVE_PRIVATE_REPOSTATS', res.body);
          callback();
        }
      });
    } else {
      callback(null, null);
    }
  };
}

export default function repoSettingsMain({actionContext, payload, done, maybeData}){
  debug('maybeData:', maybeData);
  if (_.has(maybeData, 'token')) {
    let args = {
      actionContext,
      maybeToken: maybeData.token,
      user: payload.params.user,
      splat: payload.params.splat
    };

    waterfall([
      getRepo(args),
      handleGetPrivateRepoStats(args)
    ], function(err, res){
      done();
    });
  } else {
    actionContext.dispatch('REPO_NOT_FOUND', null);
    done();
  }
}
