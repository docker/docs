'use strict';
import createStore from 'fluxible/addons/createStore';

var BillingPlansStore = createStore({
  storeName: 'CloudBillingStore',
  handlers: {
    RESET_CLOUD_BILLING_PLANS: '_clearStore',
    RECEIVE_CLOUD_BILLING_INFO: '_receiveBillingInfo',
    LOGOUT: '_clearStore'
  },
  initialize: function() {
    this.currentPlan = {};
    this.billingInfo = {};
    this.accountInfo = {};
  },
  _clearStore: function() {
    this.initialize();
    this.emitChange();
  },
  _receiveBillingInfo: function(payload) {
    this.billingInfo = payload.billingInfo;
    this.accountInfo = payload.accountInfo;
    this.currentPlan = payload.currentPlan;
    this.emitChange();
  },
  getState: function() {
    return {
      currentPlan: this.currentPlan,
      accountInfo: this.accountInfo,
      billingInfo: this.billingInfo
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.currentPlan = state.currentPlan;
    this.accountInfo = state.accountInfo;
    this.billingInfo = state.billingInfo;
  }
});

module.exports = BillingPlansStore;
