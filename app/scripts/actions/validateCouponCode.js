'use strict';
var debug = require('debug')('hub:actions:validateCouponCode');

module.exports = (actionContext, {coupon_code, plan}) => {
  if (coupon_code) {
    try {
      /**
       * This will only run where window is defined. ie: the browser
       * It throws an exception in node
       */
      window.recurly.configure(process.env.RECURLY_PUBLIC_KEY);
      /* eslint-disable new-cap */
      var pricing = window.recurly.Pricing();
      /* eslint-enable new-cap*/
      pricing
        .plan(plan, {quantity: 1})
        .coupon(coupon_code)
        .catch(function(err) {
          debug('ERROR', err);
          actionContext.dispatch('BILLING_ERRORS', {fieldErrors: {coupon_code: true}});
        })
        .done(function(price) {
          var discount = price.now.discount;
          if (discount % 1 === 0) {
            // remove the trailing decimals if its a whole number
            discount = Math.round(discount);
          }
          actionContext.dispatch('UPDATE_COUPON_VALUE', discount);
          if (discount <= 0) {
            actionContext.dispatch('BILLING_ERRORS', {fieldErrors: {coupon_code: true}});
          } else {
            actionContext.dispatch('BILLING_ERRORS', {fieldErrors: {coupon_code: false}});
          }
        });
    } catch(e) {
      debug('error', e);
      actionContext.dispatch('BILLING_ERRORS', {fieldErrors: {coupon_code: true}});
    }
    actionContext.dispatch('BILLING_INFO_UPDATE_FIELD_WITH_VALUE', {field: 'card', fieldKey: 'coupon_code', fieldValue: coupon_code});
  } else {
    actionContext.dispatch('BILLING_ERRORS', {fieldErrors: {coupon_code: false}});
    actionContext.dispatch('BILLING_INFO_UPDATE_FIELD_WITH_VALUE', {field: 'card', fieldKey: 'coupon_code', fieldValue: ''});
    actionContext.dispatch('UPDATE_COUPON_VALUE', 0);
  }
};
