'use strict';
/*global UpdateBillingInfoPayload */
import { Billing } from 'hub-js-sdk';
import _ from 'lodash';
import async from 'async';
var debug = require('debug')('hub:actions:updateBillingInformation');

function updateBillingInformation(actionContext, {JWT, user, billingInfo, accountInfo, card}, done) {
  actionContext.dispatch('BILLING_SUBMIT_START');
  var recurlyData = _.merge({}, billingInfo, card);
  let isOrg = false;
  let username = '';
  if (_.has(user, 'username')) {
    username = user.username;
  } else if (_.has(user, 'orgname')) {
    username = user.orgname;
    isOrg = true;
  }
  const account = _.merge({}, accountInfo, {user: username});

  try {
    /**
     * This will only run where window is defined. ie: the browser
     * It throws an exception in node
     */
    window.recurly.configure(process.env.RECURLY_PUBLIC_KEY);
  } catch(e) {
    debug('error', e);
  }
  window.recurly.token(recurlyData, function(recurlyErr, token) {
    if (recurlyErr) {
      // SHOULD ONLY GET HERE IF RECURLY DECLINES THEIR INFO
      debug('recurly error', recurlyErr);
      actionContext.dispatch('GET_RECURLY_ERROR', recurlyErr);
      done();
    } else {
      async.series([
        function(callback){
          Billing.updateBillingAccount(JWT, username, account, function(accountErr, accountRes){
            if (accountErr) {
              let message = 'There was an error updating your contact information. Please check your information and try again.';
              if (_.has(accountRes.body, 'detail') && _.isString(accountRes.body.detail)) {
                message = accountRes.body.detail;
              }
              actionContext.dispatch('BILLING_SUBMIT_ERROR', message);
              callback(accountErr, accountRes);
            } else if (accountRes.ok) {
              callback(null, accountRes.body);
            }
          });
        },
        function(callback) {
          let updatedInfo = {username: username, payment_token: token.id};
          Billing.updateBillingInfo(JWT, username, updatedInfo, function(billErr, billRes) {
            if (billErr) {
              debug('updateBillingInfo error', billErr);
              let message = 'There was an error submitting your billing information. Please check your information and try again.';
              if (_.has(billRes.body, 'detail') && _.isString(billRes.body.detail)) {
                message = billRes.body.detail;
              }
              actionContext.dispatch('BILLING_SUBMIT_ERROR', message);
              callback(billErr, billRes);
            } else if (billRes.ok) {
              callback(null, billRes.body);
            }
          });
        }
      ],
      function(err, res) {
        if (err) {
          done();
        } else {
          if (isOrg) {
            actionContext.history.push(`/u/${username}/dashboard/billing/`);
          } else {
            actionContext.history.push('/account/billing-plans/');
          }
          actionContext.dispatch('RECEIVE_BILLING_INFO', {billingInfo: res[1], accountInfo: res[0]});
          actionContext.dispatch('BILLING_SUBMIT_SUCCESS');
          done();
        }
      });
    }
  });
}

module.exports = updateBillingInformation;
