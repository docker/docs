'use strict';

import { parallel, waterfall } from 'async';
import sortBy from 'lodash/collection/sortBy';
import { Auth,
         Repositories as Repos,
         Billing
       } from 'hub-js-sdk';
var debug = require('debug')('hub:actions:enterpriseAttemptLogin');
import { getActivityFeed } from 'hub-js-sdk/src/Hub/SDK/Notifications';
import { getUser } from 'hub-js-sdk/src/Hub/SDK/JWT';
import {
  getNamespacesForUser
} from 'hub-js-sdk/src/Hub/SDK/Users';
import { getToken } from 'hub-js-sdk/src/Hub/SDK/Auth';
import request from 'superagent';

//Get orgs for user
function _getOrgsForCurrentUser({jwt, username, dispatch}) {
  return function(user, cb) {
    getNamespacesForUser(jwt, function(err, res) {
      if (err) {
        debug('getNamespacesForUser', err);
        cb(null);
      } else {
        // brute force the namespace reception so we can use a single action for now
        dispatch('ENTERPRISE_TRIAL_RECEIVE_ORGS', res.body.namespaces);
        dispatch('ENTERPRISE_PAID_RECEIVE_ORGS', res.body.namespaces);
        cb(null, user);
      }
    });
  };
}

function _getBillingPlans({jwt, dispatch}) {
  return function(user, cb) {
    Billing.getPlans(jwt, 'personal', (err, res) => {
      if(err) {
        debug('getPlans error', err);
        cb(null);
      } else {
        let plansList = res.body;
        let sortedPlans = sortBy(plansList, 'display_order');
        dispatch('RECEIVE_BILLING_PLANS', {
          plansList: sortedPlans
        });
        cb(null, user);
      }
    });
  };
}
function _getBillingAccount({jwt, dispatch}) {
  return function(user, cb) {
    Billing.getBillingAccount(jwt, user.username, (err, res) => {
      if(err) {
        debug('no billing account connected');
        cb(err, {});
      } else {
        dispatch('RECEIVE_BILLING_INFO', {
          accountInfo: res.body
        });
        cb(null, user);
      }
    });
  };
}

function _getBillingInfo({jwt, dispatch}) {
  return (user, cb) => {
    Billing.getBillingInfo(jwt, user.username, (err, res) => {
      if(err) {
        debug('no billing account connected');
        cb(err, {});
      } else {
        dispatch('RECEIVE_BILLING_INFO', {
          billingInfo: res.body
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
          dispatch('RECEIVE_USER', res.body); //LOGIN USING JWT.getUser
          cb(null, res.body);
        }
      });
    },
    _getOrgsForCurrentUser({jwt, username, dispatch}),
    _getBillingPlans({jwt, dispatch}),
    _getBillingAccount({jwt, dispatch}),
    _getBillingInfo({jwt, dispatch}) //Waterfalling the user through each call.
  ], function(err, user) {
    callback(err, { user });
  });
}

module.exports = function({ dispatch },
                          { username, password },
                          done) {
  dispatch('LOGIN_ATTEMPT_START');
  getToken(username,
           password,
           function(err, res) {
             if (err) {
               debug('getToken error', err);
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
               if (res.body.token) {
                 request.post('/attempt-login/')
                 .send({jwt: res.body.token})
                 .end((cookieErr, cookieRes) => {
                   handleGetUserInfo({
                     jwt: res.body.token,
                     username,
                     dispatch
                   }, function(userErr, userRes) {
                     /**
                     * CreateBillingSubscription requires having allPlans
                     * If loggedOut and JWT is populated first, it will attempt to populate form, BEFORE plans have been dispatched
                     */
                     dispatch('RECEIVE_JWT', res.body.token);
                     dispatch('LOGIN_CLEAR');
                   });
                 });
               }
             }
           });
};
