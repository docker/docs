'use strict';

import _ from 'lodash';
import async from 'async';
import { PENDING_DELETE } from 'common/enums/RepoStatus';
import {
  Repositories as Repos,
  Autobuilds
} from 'hub-js-sdk';

export default function repo({actionContext, payload, done, maybeData}){
  var token;
  if (_.has(maybeData, 'token')) {
    token = maybeData.token;
  } else {
    token = '';
  }

  var namespace = payload.params.user;
  if(payload.params.user === '_') {
    namespace = 'library';
  }
  var repoShortName = namespace + '/' + payload.params.splat;

  //1. Get repository details
  //2. If successful, set valid repo and pass it along for the next set of calls
  //3. If not valid, set repo to be not valid and dispatch repo not found
  var _getRepo = function(cb) {
    Repos.getRepo(token, repoShortName, function(err, res) {
      var repoInfo = {};
      let status;
      if (res && res.body) {
        status = res.body.status;
      }

      if (err || status === PENDING_DELETE) {
        repoInfo.error = err;
        repoInfo.isValid = false;
        cb(null, repoInfo);
        actionContext.dispatch('REPO_NOT_FOUND', err);
      } else {
        repoInfo.info = res.body;
        actionContext.dispatch('RECEIVE_REPOSITORY', res.body);
        repoInfo.isValid = true;
        cb(null, repoInfo);
      }
    });
  };

  var _getRepoDetails = function(repoInfo, cb) {
    if (repoInfo.isValid) {
      async.parallel([
          function(callback) {
            Repos.getCommentsForRepo(token, repoShortName, function(err, res) {
              if (err) {
                callback(err);
              } else {
                callback(null, res.body);
              }
            }, 1);
          }, function(callback) {
            if (repoInfo.info.is_automated) {
              Autobuilds.getAutomatedBuildSettings(token, namespace, payload.params.splat, function (err, res) {
                if (err) {
                  actionContext.dispatch('AUTOBUILD_REPO_NOT_FOUND');
                  callback(null);
                } else {
                  callback(null, res.body);
                }
              });
            } else {
              callback(null);
            }
          }
        ],
        function(error, results) {
          if (error) {
            cb(error);
          } else {
            cb(null, results);
          }
        });
    } else {
      cb('repo not found');
    }
  };

  async.waterfall([
    _getRepo,
    _getRepoDetails
  ], function(e, finalResults) {
    if (finalResults) {
      actionContext.dispatch('RECEIVE_REPO_COMMENTS', finalResults[0]);
      if (finalResults[1]) {
        actionContext.dispatch('RECEIVE_AUTOBUILD_SETTINGS', finalResults[1]);
      }
    }
    done();
  });
}
