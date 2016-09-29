'use strict';

const debug = require('debug')('navigate::billingPlans');
import async from 'async';
import _ from 'lodash';

import {
  Billing,
  Users
  } from 'hub-js-sdk';

//GETs current HUB subscription
function _getBillingSubscriptions({token, user}) {
  return (callback) => {
    Billing.getBillingSubscriptions(token, user.username, (err, res) => {
      if(err) {
        callback(null, {});
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
      if(err) {
        debug('NO BILLING ACCOUNT CONNECTED');
        callback(null, {});
      } else {
        debug('GET BILLING ACCOUNT', res.body);
        callback(null, res.body);
      }
    });
  };
}

function _getBillingInfo({token, user}) {
  return (callback) => {
    Billing.getBillingInfo(token, user.username, (err, res) => {
      if(err) {
        debug('NO BILLING ACCOUNT CONNECTED');
        callback(null, {});
      } else {
        debug('GET BILLING INFO', res.body);
        callback(null, res.body);
      }
    });
  };
}
/**
 * =============================================
 *  ^ GET SUBSCRIPTIONS/ACCOUNTINFO/BILLINGINFO
 *  Must re-get on change of namespace
 * =============================================
 */

//GET namespaces for user
function _getNamespaces({token}) {
  return (callback) => {
    Users.getNamespacesForUser(token, function(err, res) {
      if (err) {
        callback(null, {});
      } else {
        callback(null, res.body.namespaces);
      }
    });
  };
}

export default function billingPlans({actionContext, payload, done, maybeData}){
  if (maybeData.token && maybeData.user) {
    const {token, user} = maybeData;
    actionContext.dispatch('RESET_CLOUD_BILLING_PLANS');
    async.parallel({
      userPlan: _getBillingSubscriptions({token, user}),
      accountInfo: _getBillingAccount({token, user}),
      billingInfo: _getBillingInfo({token, user}),
      namespaces: _getNamespaces({token})
    }, function(err, results){
      const { userPlan, accountInfo, billingInfo, namespaces } = results;
      debug('CLOUD BILLING PLANS', results);
      actionContext.dispatch('RECEIVE_CLOUD_BILLING_INFO', {
        billingInfo: billingInfo,
        accountInfo: accountInfo,
        currentPlan: userPlan
      });
      actionContext.dispatch('ENTERPRISE_PAID_RECEIVE_ORGS', namespaces);
      const values = _.merge({}, billingInfo, {
        account_first: accountInfo.first_name,
        account_last: accountInfo.last_name,
        company_name: accountInfo.company_name,
        email: accountInfo.email
      });
      debug('INITIALIZE ENTERPRISE BILLING FORM: ', values);
      // Need to differentiate btw billingInfo/accountInfo first/last names
      actionContext.dispatch('ENTERPRISE_PAID_POPULATE_FORM', values);
      return done();
    });
  } else {
    done();
  }
}
