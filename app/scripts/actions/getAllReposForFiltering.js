'use strict';
var debug = require('debug')('hub:actions:getAllReposForFiltering');
import async from 'async';
import isArray from 'lodash/lang/isArray';
import {
  Repositories as Repos
} from 'hub-js-sdk';

export default function getAllReposForFiltering(actionContext, {jwt, user}) {

    actionContext.dispatch('DASHBOARD_REPOS_STORE_ATTEMPTING_GET_ALL_REPOS');
    var _getAllRepos = function (cb) {
      Repos.getAllReposForUser(jwt, user, function (err, res) {
        if (err) {
          cb(err);
        } else {
          cb(null, isArray(res.body) && res.body.length);
        }
      });
    };

    //Get repos for user or org
    var _getReposForUserOrOrg = function (pageSize, cb) {
      Repos.getReposForUser(jwt, user, function (err, res) {
        if (err) {
          actionContext.dispatch('ERROR_RECEIVING_REPOS', err);
          cb();
        } else {
          actionContext.dispatch('DASHBOARD_REPOS_STORE_RECEIVE_ALL_REPOS_SUCCESS', res.body);
          cb();
        }
      }, 1, pageSize);
    };

    async.waterfall([
      _getAllRepos,
      _getReposForUserOrOrg
    ], function(error, response) {
    });
}
