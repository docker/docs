'use strict';

import { parallel, waterfall } from 'async';
import { Auth,
         Repositories as Repos
       } from 'hub-js-sdk';
var debug = require('debug')('hub:actions:attemptLogin');
import { getActivityFeed } from 'hub-js-sdk/src/Hub/SDK/Notifications';
import { getUser } from 'hub-js-sdk/src/Hub/SDK/JWT';
import {
  getOrgsForUser,
  getUserSettings
} from 'hub-js-sdk/src/Hub/SDK/Users';
import { getToken } from 'hub-js-sdk/src/Hub/SDK/Auth';
import request from 'superagent';

function handleGetRepos({jwt, username, dispatch}) {
  return function(callback) {
    Repos.getReposForUser(jwt, username, function(err, res) {
      if (err) {
        callback(null, null);
      } else {
        dispatch('RECEIVE_REPOS', res.body);
      }
    }, 1); //get page 1
  };
}

function handleGetPrivateRepoStats({jwt, username, dispatch}) {
  return function(callback) {
    getUserSettings(jwt, username, function(err, res) {
      if (res.ok) {
        dispatch('RECEIVE_PRIVATE_REPOSTATS', res.body);
      }
      callback(null, null);
    });
  };
}

//Get orgs for user
function _getOrgsForCurrentUser({jwt, username, dispatch}) {
  return function(user, cb) {
    getOrgsForUser(jwt, function(err, res) {
      if (err) {
        debug('error', err);
        cb(null);
      } else {
        dispatch('RECEIVE_DASHBOARD_NAMESPACES', {
          orgs: res.body,
          user: user.username
        });
        cb(null, user);
      }
    });
  };
}

function handleGetUserInfo({jwt, username, dispatch}, callback) {
  waterfall([
    function(cb){
      getUser(jwt, function(err, res) {
        if (err) {
          cb(err, {});
        } else {
          dispatch('RECEIVE_USER', res.body);
          cb(null, res.body);
        }
      });
    },
    _getOrgsForCurrentUser({jwt, username, dispatch})
  ], function(err, user) {
    callback(err, { user });
  });
}

function fetchDataForDashboard({
  jwt, username, dispatch, done
}) {
  dispatch('RECEIVE_JWT', jwt);
  dispatch('DASHBOARD_REPOS_STORE_ATTEMPTING_GET_REPOS');
  parallel([
    handleGetRepos({jwt, username, dispatch}),
    handleGetPrivateRepoStats({jwt, username, dispatch})
  ], function(err, results) {
    if (err) {
      debug('error', err);
    }
    return done();
  });
}

module.exports = function({ dispatch, history },
                          { username, password },
                          done) {
  dispatch('LOGIN_ATTEMPT_START');
  getToken(username,
           password,
           function(err, res) {
             if (err) {
               debug('error', err);
               if (res.unauthorized) {
                 if(res.body && res.body.detail) {
                   /**
                    * This can happen if the user has not verified their email
                    */
                   dispatch('LOGIN_UNAUTHORIZED_DETAIL', res.body);
                 } else {
                   dispatch('LOGIN_UNAUTHORIZED');
                 }
               } else if (res.badRequest){
                 try {
                   dispatch('LOGIN_BAD_REQUEST', JSON.parse(res.text));
                 } catch (error) {
                   dispatch('LOGIN_ERROR');
                 }
               } else {
                 // unhandled login error
                 dispatch('LOGIN_ERROR');
               }
             } else {
               debug('got token');
               if (res.body.token) {
                 request.post('/attempt-login/')
                 .send({jwt: res.body.token})
                 .end((cookieErr, cookieRes) => {
                   handleGetUserInfo({
                     jwt: res.body.token,
                     username,
                     dispatch
                   }, function(userErr, userRes) {
                     dispatch('LOGIN_CLEAR');
                     dispatch('CURRENT_USER_CONTEXT', {
                       username: userRes.user.username
                     });
                     history.pushState(null, '/');
                     fetchDataForDashboard({
                       jwt: res.body.token,
                       username: userRes.user.username,
                       dispatch,
                       done
                     });
                   });
                 });
               }
             }
           });
};
