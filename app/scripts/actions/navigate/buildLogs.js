'use strict';

import has from 'lodash/object/has';
import { parallel, waterfall } from 'async';
import request from 'superagent';
import { PENDING_DELETE } from 'common/enums/RepoStatus';
import {
  Repositories as Repos,
  Autobuilds
} from 'hub-js-sdk';
const debug = require('debug')('buildLogs');

export default function buildLogs({actionContext, payload, done, maybeData}) {
  var token;
  if (has(maybeData, 'token')) {
    token = maybeData.token;
  } else {
    token = '';
  }

  var namespace = payload.params.user;
  if (payload.params.user === '_') {
    namespace = 'library';
  }
  var repoShortName = namespace + '/' + payload.params.splat;
  const build_code = payload.params.build_code;

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
            Repos.getCommentsForRepo(token, repoShortName, function (err, res) {
              if (err) {
                callback(err);
              } else {
                callback(null, res.body);
              }
            });
          }, function (callback) {
            request.get(process.env.REGISTRY_API_BASE_URL + '/v2/repositories/' + repoShortName + '/buildhistory/' + build_code + '/')
                   .accept('application/json')
                   .set('Authorization', 'JWT ' + token)
                   .end((err, res) => {
                     if(err) {
                       callback();
                     } else {
                       debug('BUILD LOGS RECEIVE', res.body);
                       actionContext.dispatch('BUILD_LOGS_RECEIVE', res.body);
                       callback();
                     }
                   });
          }, function (callback) {
            if (repoDetail.is_automated) {
              Autobuilds.getAutomatedBuildSettings(token, namespace, payload.params.splat, function (err, res) {
                if (err) {
                  actionContext.dispatch('AUTOBUILD_REPO_NOT_FOUND');
                  callback(err);
                } else {
                  actionContext.dispatch('RECEIVE_AUTOBUILD_SETTINGS', res.body);
                  callback();
                }
              });
            } else {
              callback();
            }
          }
        ],
        function (error, results) {
          if (error) {
            cb();
          } else {
            actionContext.dispatch('RECEIVE_REPO_COMMENTS', results[0]);
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
