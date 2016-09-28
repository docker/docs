'use strict';
/**
 * CLONE OF updateSubscriptionPlanOrPackage.js
 * with Different dispatches
 */
import { Billing } from 'hub-js-sdk';
import has from 'lodash/object/has';
var debug = require('debug')('hub:actions:updateCloudBillingSubscription');

function updateBillingSubscriptions(actionContext, {JWT, username, subscription_uuid, package_code, coupon_code, isOrg}, done) {
  actionContext.dispatch('ENTERPRISE_PAID_ATTEMPT_START');
  var data = {
    package: package_code,
    username,
    coupon_code
  };
  Billing.updateBillingSubscriptions(JWT, subscription_uuid, username, data, function(err, res) {
    if (err) {
      debug(err);
      // If fails - There's nothing else we can do. Facepalm
      actionContext.dispatch('ENTERPRISE_PAID_BAD_REQUEST', res.body.detail);
      done();
    } else if (res.ok) {
      actionContext.dispatch('ENTERPRISE_PAID_SUCCESS');
      if (isOrg) {
        actionContext.history.push(`/u/${username}/dashboard/billing/`);
      } else {
        actionContext.history.push('/account/billing-plans/');
      }
      done();
    }
  });
}

module.exports = updateBillingSubscriptions;
