'use strict';
import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('PlansStore');

var PlansStore = createStore({
  storeName: 'PlansStore',
  handlers: {
    RECEIVE_BILLING_PLANS: '_receivePlans',
    RECEIVE_BILLING_INFO: '_receiveBilling',
    RECEIVE_BILLING_SUBSCRIPTION: '_receiveBillingSubscription',
    RESET_CURRENT_PLAN: '_resetCurrentPlan'
  },
  initialize: function() {
    this.plansList = [];
    this.currentPlan = {
      plan: '',
      package: '',
      subscription_uuid: '',
      state: '',
      add_ons: []
    };
  },
  _clearPlan: function() {
    this.currentPlan = {
      plan: '',
      package: '',
      subscription_uuid: '',
      state: '',
      add_ons: []
    };
  },
  _receiveBilling: function(payload){
    debug('RECEIVE BILLING: ', payload);
    this._clearPlan();
    if (payload.currentPlan) {
      this.currentPlan = payload.currentPlan;
    }
    this.emitChange();
  },
  _receiveBillingSubscription: function(payload) {
    this._clearPlan();
    if (payload.currentPlan) {
      this.currentPlan = payload.currentPlan;
    }
    this.emitChange();
  },
  _resetCurrentPlan: function(payload) {
    this._clearPlan();
    if (payload.currentPlan) {
      this.currentPlan = payload.currentPlan;
    }
    this.emitChange();
  },
  _receivePlans: function(payload) {
    debug(payload);
    this.plansList = payload.plansList;
    this.emitChange();
    },
  getState() {
    return {
      plansList: this.plansList,
      currentPlan: this.currentPlan
    };
  },
  dehydrate() {
    return this.getState();
  },
  rehydrate(state) {
    this.plansList = state.plansList;
    this.currentPlan = state.currentPlan;
  }
});

module.exports = PlansStore;
