'use strict';
import createStore from 'fluxible/addons/createStore';
import _ from 'lodash';

var CloudCouponStore = createStore({
  storeName: 'CloudCouponStore',
  handlers: {
    RECEIVE_BILLING_INFO: '_clearStore',
    CLEAR_CLOUD_COUPON: '_clearStore',
    UPDATE_COUPON_VALUE: '_updateDiscountValue',
    BILLING_INFO_UPDATE_FIELD_WITH_VALUE: '_updateCouponCode',
    BILLING_ERRORS: '_updateErrors'
  },
  initialize: function() {
    this.couponCode = '';
    this.discountValue = 0;
    this.hasError = false;
  },
  _clearStore: function() {
    this.initialize();
    this.emitChange();
  },
  _updateCouponCode: function({field, fieldKey, fieldValue}) {
    this.couponCode = fieldKey === 'coupon_code' ? fieldValue : this.couponValue;
    this.emitChange();
  },
  _updateDiscountValue: function(discount) {
    this.discountValue = discount;
    this.emitChange();
  },
  _updateErrors: function({fieldErrors}) {
    if (_.has(fieldErrors, 'coupon_code')) {
      this.hasError = fieldErrors.coupon_code;
    }
    this.emitChange();
  },
  getState: function() {
    return {
      couponCode: this.couponCode,
      discountValue: this.discountValue,
      hasError: this.hasError
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.couponCode = state.couponCode;
    this.discountValue = state.discountValue;
    this.hasError = state.hasError;
  }
});

module.exports = CloudCouponStore;
