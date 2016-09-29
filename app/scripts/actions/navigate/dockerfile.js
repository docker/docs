'use strict';

import has from 'lodash/object/has';
import { parallel, waterfall } from 'async';
import { PENDING_DELETE } from 'common/enums/RepoStatus';
import {
  Repositories as Repos,
  Autobuilds
} from 'hub-js-sdk';

export default function dockerfile({actionContext, payload, done, maybeData}){
  var token;
  if (has(maybeData, 'token')) {
    token = maybeData.token;
  } else {
    token = '';
  }

  var namespace = payload.params.user;
  if(payload.params.user === '_') {
    namespace = 'library';
  }
  var repoShortName = namespace + '/' + payload.params.splat;

  var _getRepo = function (callback) {
    Repos.getRepo(token, repoShortName, function (err, res) {
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
  };

  var _getRest = function (repoDetail, cb) {
    if (repoDetail) {
      parallel([
          function (callback) {
            Repos.getDockerfile(token, repoShortName, function (err, res) {
              if (err) {
                callback();
              } else {
                actionContext.dispatch('RECEIVE_DOCKERFILE_FOR_REPOSITORY', res.body);
                callback();
              }
            });
          }, function(callback) {
            if (repoDetail.is_automated) {
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
        function (error, results) {
          if (error) {
            cb();
          } else {
            if (results[1]) {
              actionContext.dispatch('RECEIVE_AUTOBUILD_SETTINGS', results[1]);
            }
            cb();
          }
        });
    } else {
      cb();
    }
  };

  waterfall([
    _getRepo,
    _getRest
  ], function(err, res) {
    done();
  });
}
