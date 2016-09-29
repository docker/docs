'use strict';
import createStore from 'fluxible/addons/createStore';
import _ from 'lodash';

var BillingPlansStore = createStore({
  storeName: 'BillingPlansStore',
  handlers: {
    RESET_BILLING_PLANS: '_clearStore',
    RECEIVE_BILLING_INFO: '_receiveBillingInfo',
    RECEIVE_BILLING_SUBSCRIPTION: '_receiveBillingSubscription',
    RECEIVE_INVOICES: '_receiveInvoices',
    RESET_CURRENT_PLAN: '_resetCurrentPlan',
    UPDATE_PLAN_START: '_updatePlanStart',
    UPDATE_PLAN_ERROR: '_updatePlanErr',
    DELETE_SUBSCRIPTION_ERR: '_updatePlanErr',
    DELETE_SUBSCRIPTION_SUCCESS: '_unsubscribeComplete',
    UNSUBSCRIBE_SUBSCRIPTION: '_unsubscribe', // UNUSED - deprecated
    UNSUBSCRIBE_PACKAGE: '_unsubscribePackage', // UNUSED - deprecated
    UNSUBSCRIBE_PLAN: '_unsubscribePlan', // UNUSED - deprecated
    LOGOUT: '_clearStore'
  },
  initialize: function() {
    this.currentPlan = {};
    this.accountInfo = {
      account_code: '',
      username: '',
      email: '',
      first_name: '',
      last_name: '',
      company_name: '',
      hasError: false,
      newBilling: true
    };
    this.billingInfo = {
      first_name: '',
      last_name: '',
      address1: '',
      address2: '',
      country: '',
      state: '',
      zip: '',
      city: '',
      last_four: '',
      card_type: '',
      month: '',
      year: '',
      newBilling: true
    };
    this.invoices = [];
    this.plansError = '';
    this.unsubscribing = '';
    this.updatePlan = '';
  },
  _clearStore: function() {
    this.initialize();
  },
  _receiveBillingInfo: function(payload) {
    this.initialize();
    _.merge(this.billingInfo, payload.billingInfo);
    _.merge(this.accountInfo, payload.accountInfo);
    _.merge(this.currentPlan, payload.currentPlan);
    this.emitChange();
  },
  _receiveBillingSubscription: function(payload) {
    _.merge(this.currentPlan, payload.currentPlan);
    this.updatePlan = '';
    this.emitChange();
  },
  _receiveInvoices: function(payload) {
    this.invoices = payload.invoices;
    this.emitChange();
  },
  _resetCurrentPlan: function(payload) {
    this.currentPlan = payload.currentPlan;
    this.emitChange();
  },
  _unsubscribe: function() { // UNUSED - deprecated
    this.unsubscribing = 'subscription';
    this.emitChange();
  },
  _unsubscribePackage: function() { // UNUSED - deprecated
    this.unsubscribing = 'package';
    this.emitChange();
  },
  _unsubscribePlan: function() { // UNUSED - deprecated
    this.unsubscribing = 'plan';
    this.emitChange();
  },
  _unsubscribeComplete: function() { // UNUSED - deprecated
    this.unsubscribing = '';
    this.emitChange();
  },
  _updatePlanStart: function(payload) {
    this.updatePlan = payload;
    this.emitChange();
  },
  _updatePlanErr: function(payload) {
    this.unsubscribing = '';
    this.updatePlan = '';
    this.plansError = payload;
    this.emitChange();
  },
  getState: function() {
    return {
      currentPlan: this.currentPlan,
      accountInfo: this.accountInfo,
      billingInfo: this.billingInfo,
      invoices: this.invoices,
      plansError: this.plansError,
      unsubscribing: this.unsubscribing,
      updatePlan: this.updatePlan
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.currentPlan = state.currentPlan;
    this.accountInfo = state.accountInfo;
    this.billingInfo = state.billingInfo;
    this.invoices = state.invoices;
    this.plansError = state.plansError;
    this.unsubscribing = state.unsubscribing;
    this.updatePlan = state.updatePlan;
  }
});

module.exports = BillingPlansStore;
