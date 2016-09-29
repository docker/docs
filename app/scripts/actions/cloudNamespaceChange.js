'use strict';

const debug = require('debug')('hub:actions:cloudNamespaceChange');
import async from 'async';
import _ from 'lodash';

import {
  Billing,
  Users
  } from 'hub-js-sdk';

//GETs current HUB subscription
function _getBillingSubscriptions({JWT, namespace}) {
  return (callback) => {
    Billing.getBillingSubscriptions(JWT, namespace, (err, res) => {
      if(err) {
        callback(null, {});
      } else {
        let subscriptions = res.body;

        // the assumption is that there is at most one subscription right now
        let subscription = _.head(subscriptions);
        callback(null, subscription);
      }
    });
  };
}

function _getBillingAccount({JWT, namespace}) {
  return (callback) => {
    Billing.getBillingAccount(JWT, namespace, (err, res) => {
      if(err) {
        debug('getBillingAccount', 'no billing account connected');
        callback(null, {});
      } else {
        callback(null, res.body);
      }
    });
  };
}

function _getBillingInfo({JWT, namespace}) {
  return (callback) => {
    Billing.getBillingInfo(JWT, namespace, (err, res) => {
      if(err) {
        debug('getBillingInfo', 'no billing account connected');
        callback(null, {});
      } else {
        callback(null, res.body);
      }
    });
  };
}

export default function billingPlans(actionContext, {JWT, namespace}, done){
  actionContext.dispatch('RESET_CLOUD_BILLING_PLANS');
  async.parallel({
    userPlan: _getBillingSubscriptions({JWT, namespace}),
    accountInfo: _getBillingAccount({JWT, namespace}),
    billingInfo: _getBillingInfo({JWT, namespace})
  }, function(err, results){
    let { userPlan, accountInfo, billingInfo } = results;
    actionContext.dispatch('RECEIVE_CLOUD_BILLING_INFO', {
      billingInfo: billingInfo,
      accountInfo: accountInfo,
      currentPlan: userPlan
    });
    let values = _.merge({}, billingInfo, {
      account_first: accountInfo.first_name,
      account_last: accountInfo.last_name,
      company_name: accountInfo.company_name,
      email: accountInfo.email
    });
    // TODO: Need to differentiate btw billingInfo/accountInfo first/last names
    actionContext.dispatch('ENTERPRISE_PAID_POPULATE_FORM', values);
    return done();
  });
}
