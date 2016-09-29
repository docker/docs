'use strict';
var debug = require('debug')('navigate::userStars');
import async from 'async';
import _ from 'lodash';
import {
  Repositories as Repos,
  Users
} from 'hub-js-sdk';

export default function userStars({actionContext, payload, done, maybeData}){
  // This works without a jwt
  var token = null;

  if (_.has(maybeData, 'token')) {
    token = maybeData.token;
  }

  //Get repos for user or org
  var _getStarredReposForUser = function(cb) {
    Repos.getStarredReposForUser(token, payload.params.user, function(err, res) {
      if (err) {
        cb();
      } else {
        actionContext.dispatch('RECEIVE_PROFILE_STARRED_REPOS', res.body);
        cb();
      }
    }, payload.location.query.page);
  };

  var _getUserByName = function(cb) {
    Users.getUser(token, payload.params.user, function(err, res) {
      if (err) {
        if(res && res.notFound) {
          actionContext.dispatch('USER_PROFILE_404');
          cb(err, null);
        } else {
          debug(err);
          cb();
        }
      } else {
        actionContext.dispatch('RECEIVE_PROFILE_USER', res.body);
        cb();
      }
    });
  };

  //Get user by name
  async.parallel([
    _getStarredReposForUser,
    _getUserByName
  ], function(err, results) {
    if (err) {
      debug(err);
    }
    return done();
  });
}
