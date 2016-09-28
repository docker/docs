'use strict';

const debug = require('debug')('navigate::billingPlans');
import async from 'async';
import _ from 'lodash';
import {
  BILLFORWARD_ACCOUNT_ID
} from 'stores/common/Constants.js';

import {
  Billing
} from 'hub-js-sdk';

function _getPersonalPlans({token}) {
  return (callback) => {
    Billing.getPlans(token, 'personal', (err, res) => {
      if(err) {
        debug(err);
        callback(null, []);
      } else {
        let plansList = res.body;
        let sortedPlans = _.sortBy(plansList, 'display_order');
        callback(null, sortedPlans);
      }
    });
  };
}

function _getBillingSubscriptions({token, user}) {
  return (callback) => {
    Billing.getBillingSubscriptions(token, user.username, (err, res) => {
      if(!res || res.status === 500) {
        debug('SERVER ISSUE: getting subscriptions');
        callback('SERVER ISSUE');
      } else if (err) {
        callback(null, res.body);
      } else {
        debug('GET USER PLANS', res.body);
        let subscriptions = res.body;

        // the assumption is that there is at most one subscription right now
        let subscription = _.head(subscriptions);
        callback(err, subscription);
      }
    });
  };
}

function _getBillingAccount({token, user}) {
  return (callback) => {
    Billing.getBillingAccount(token, user.username, (err, res) => {
      if(!res || res.status === 500) {
        debug('SERVER ISSUE: getting billing account info');
        // NOTE: If this is a brand new billing account - set newBilling to true
        callback(null, {account: { newBilling: true }});
      } else if (err) {
        debug('NO BILLING ACCOUNT CONNECTED');
        // NOTE: If this is a brand new billing account - set newBilling to true
        callback(null, {account: { newBilling: true }});
      } else {
        debug('GET BILLING ACCOUNT', res.body);
        const account = { ...res.body, newBilling: false };
        callback(err, { account, billforwardId: res.header[BILLFORWARD_ACCOUNT_ID] });
      }
    });
  };
}

function _getBillingInfo({token, user}) {
  return (callback) => {
    Billing.getBillingInfo(token, user.username, (err, res) => {
      if(!res || res.status === 500) {
        debug('SERVER ISSUE: getting billing info');
        // NOTE: If this is a brand new billing profile - set newBilling to true
        callback(null, { newBilling: true });
      } else if (err) {
        debug('NO BILLING INFO CONNECTED');
        // NOTE: If this is a brand new billing profile - set newBilling to true
        callback(null, { newBilling: true });
      } else {
        debug('GET BILLING INFO', res.body);
        callback(err, { ...res.body, newBilling: false });
      }
    });
  };
}

function _getBillingInvoices({token, user}) {
  return (callback) => {
    Billing.getBillingInvoices(token, user.username, (err, res) => {
      if(err) {
        debug('NO BILLING ACCOUNT CONNECTED');
        callback(null, []);
      } else {
        debug('GET BILLING INVOICES', res.body);
        callback(err, res.body);
      }
    });
  };
}

export default function billingPlans({actionContext, payload, done, maybeData}){
  if (maybeData.token && maybeData.user) {
    let {token, user} = maybeData;
    actionContext.dispatch('RESET_BILLING_PLANS');
    /*
    NOTE:
    None of the functions should pass an error into the callback or else it will
    kill the rest of the calls and return before we're able to fetch all the data
    */
    async.parallel({
      allPlans: _getPersonalPlans({token}),
      userPlan: _getBillingSubscriptions({token, user}),
      accountInfo: _getBillingAccount({token, user}),
      billingInfo: _getBillingInfo({token, user}),
      invoiceList: _getBillingInvoices({token, user})
    }, function(err, results){
      const {
        allPlans,
        userPlan,
        accountInfo,
        billingInfo,
        invoiceList
      } = results;
      debug('BILLING PLANS', results);
      // IF AN ACCOUNT HAS BEEN MIGRATED - ITS ACCOUNT INFO WILL HAVE PARAMETER 'payment_gateway' === 'stripe'
      const gateway = accountInfo && accountInfo.account && accountInfo.account.payment_gateway;
      const billforwardId = accountInfo && accountInfo.billforwardId;
      actionContext.dispatch('RECEIVE_BILLING_PLANS', {
        plansList: allPlans
      });
      actionContext.dispatch('RECEIVE_BILLING_INFO', {
        billingInfo: billingInfo,
        accountInfo: accountInfo && accountInfo.account,
        currentPlan: userPlan,
        gateway,
        billforwardId
      });
      actionContext.dispatch('RECEIVE_INVOICES', {invoices: invoiceList});
      return done();
    });
  } else {
    done();
  }
}
