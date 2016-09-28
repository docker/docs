'use strict';
/*global UpdateBillingInfoPayload */
import { Billing } from 'hub-js-sdk';
import _ from 'lodash';
const CLOUD_METERED = 'cloud_metered';
var debug = require('debug')('hub:actions:updateSubscriptionPlanOrPackage');

function updateBillingSubscriptions(actionContext, {
    coupon_code,
    JWT,
    package_code,
    plan_code,
    username,
    subscription_uuid,
    add_ons
  }, done) {
  actionContext.dispatch('UPDATE_PLAN_START', plan_code);
  var data = {
    username,
    coupon_code,
    add_ons
  };

  // we are either updating a plan or a package
  if (plan_code) {
    data.plan = plan_code;
  }
  if (package_code) {
    data.package = package_code;
  }

  // Explicitly remove all add ons (such as nautilus) when downgrading to
  // a free account
  if (plan_code === CLOUD_METERED) {
    data.add_ons = [];
  }

  if (!plan_code || plan_code === CLOUD_METERED) {
    actionContext.dispatch('UNSUBSCRIBE_PLAN');
  } else if (!package_code) {
    actionContext.dispatch('UNSUBSCRIBE_PACKAGE');
  }
  Billing.updateBillingSubscriptions(JWT, subscription_uuid, username, data, function(err, res) {
    if (err) {
      debug(err);
      const detail = err.response && err.response.body && err.response.body.detail;
      const error = detail || 'Could not reach server.';
      actionContext.dispatch('UPDATE_PLAN_ERROR', error);
      done();
    } else if (res.ok) {
      Billing.getBillingSubscriptions(JWT, username, function(getErr, getRes) {
        if (getErr) {
          debug(getErr);
        } else if (getRes.ok) {
          let subscriptions = getRes.body;
          let subscription = _.head(subscriptions);
          // This will update your subscription in the store - while leaving the billing information the same
          actionContext.dispatch('RECEIVE_BILLING_SUBSCRIPTION', {currentPlan: subscription});
        }
        done();
      });
    }
  });
}

module.exports = updateBillingSubscriptions;
