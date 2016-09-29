'use strict';

const debug = require('debug')('navigate::OrgBilling');
import async from 'async';
import _ from 'lodash';
import {
  BILLFORWARD_ACCOUNT_ID
} from 'stores/common/Constants.js';

import {
  Billing,
  Users,
  Orgs
  } from 'hub-js-sdk';

function _getOrgPlans(token) {
  return (callback) => {
    Billing.getPlans(token, 'personal', (err, res) => {
      if(err) {
        debug(err);
        callback(null, []);
      } else {
        let plansList = res.body;
        let sortedPlans = _.sortBy(plansList, 'display_order');
        debug(sortedPlans);
        callback(null, sortedPlans);
      }
    });
  };
}

function _getBillingSubscriptions(token, username) {
  return (callback) => {
    Billing.getBillingSubscriptions(token, username, (err, res) => {
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

function _getBillingAccount(token, username) {
  return (callback) => {
    Billing.getBillingAccount(token, username, (err, res) => {
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

function _getBillingInfo(token, username) {
  return (callback) => {
    Billing.getBillingInfo(token, username, (err, res) => {
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

function _getBillingInvoices(token, username) {
  return (callback) => {
    Billing.getBillingInvoices(token, username, (err, res) => {
      if(err) {
        debug('NO BILLING ACCOUNT INVOICES');
        callback(null, []);
      } else {
        debug('GET BILLING INVOICES', res.body);
        callback(err, res.body);
      }
    });
  };
}

//Get orgs for user
function _getOrgsForCurrentUser(token) {
 return (callback) => {
  Users.getOrgsForUser(token, function(err, res) {
    if (err) {
      debug(err);
      callback();
    } else {
      debug('GET Orgs for current User', res.body);
      callback(null, res.body);
    }
  });
 };
}

//Get org's settings since we are in the organization's home
function _getOrgSettings(token, username) {
  return (callback) => {
    Orgs.getOrg(token, username, function (err, res) {
      if (err) {
        callback();
      } else {
        debug('GET Orgs Settings', res.body);
        callback(null, res.body);
      }
    });
  };
}

//Get user settings for private repo stats
var _getOrgPrivateRepoSettings = function(ac, token, username) {
  return (cb) => {
    Orgs.getOrgSettings(token, username, function (err, res) {
      if (err) {
        debug(err);
        ac.dispatch('PRIVATE_REPOSTATS_NO_PERMISSIONS', err);
        cb(null, err);
      } else {
        ac.dispatch('RECEIVE_PRIVATE_REPOSTATS', res.body);
        cb(null, res.body);
      }
    });
  };
};

var _getNamespaces = function(actionContext, token, profileuser) {
  return (callback) => {
    Users.getNamespacesForUser(token, function(err, res) {
      if (err) {
        callback();
      } else {
        callback(null, res.body);
        actionContext.dispatch('CREATE_REPO_RECEIVE_NAMESPACES', {
          namespaces: res.body,
          selectedNamespace: profileuser
        });
      }
    });
  };
};

export default function billingOrgPlans({actionContext, payload, done, maybeData}){
  if (maybeData.token && maybeData.user) {
    let {token, user} = maybeData;
    let username = payload.params.user;
    actionContext.dispatch('RESET_BILLING_PLANS');
    /*
    NOTE:
    None of the functions should pass an error into the callback or else it will
    kill the rest of the calls and return before we're able to fetch all the data
    */
    async.parallel({
      allPlans: _getOrgPlans(token),
      userPlan: _getBillingSubscriptions(token, username),
      accountInfo: _getBillingAccount(token, username),
      billingInfo: _getBillingInfo(token, username),
      invoiceList: _getBillingInvoices(token, username),
      getOrgs: _getOrgsForCurrentUser(token),
      getOrgSettings: _getOrgSettings(token, username),
      getPrivateRepoStats: _getOrgPrivateRepoSettings(actionContext, token, username),
      getNamespaces: _getNamespaces(actionContext, token, user.username)
    }, function(err, results){
      const {
        allPlans,
        userPlan,
        accountInfo,
        billingInfo,
        invoiceList,
        getOrgs,
        getOrgSettings
      } = results;
      debug('BILLING PLANS', results);
      const gateway = accountInfo && accountInfo.account && accountInfo.account.payment_gateway;
      const billforwardId = accountInfo && accountInfo.billforwardId;

      actionContext.dispatch('CURRENT_USER_CONTEXT', { username: username });
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
      actionContext.dispatch('RECEIVE_DASHBOARD_NAMESPACES', {
        orgs: getOrgs,
        user: user.username
      });
      actionContext.dispatch('RECEIVE_ORGANIZATION', getOrgSettings);
      return done();
    });
  } else {
    done();
  }
}
