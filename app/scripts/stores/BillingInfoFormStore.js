'use strict';
import createStore from 'fluxible/addons/createStore';
import _ from 'lodash';
import { STATUS } from './billingformstore/Constants';
var debug = require('debug')('STORE::BillingInfoFormStore');

var BillingInfoFormStore = createStore({
  storeName: 'BillingInfoFormStore',
  handlers: {
    BILLING_ACCOUNT_EXISTS: '_accountExists',
    BILLING_INFO_EXISTS: '_billingInfoExists',
    BILLING_ERRORS: '_updateErrors',
    BILLING_INFO_UPDATE_FIELD_WITH_VALUE: '_updateBillingInfoForm',
    BILLING_SUBMIT_ERROR: '_submitErrors',
    BILLING_SUBMIT_START: '_submitStart',
    BILLING_SUBMIT_SUCCESS: '_submitSuccess',
    CLEAR_BILLING_FORM: '_clearBillingForm',
    GET_RECURLY_ERROR: '_updateRecurlyErrors',
    LOGOUT: '_clearStore',
    RECEIVE_BILLING_INFO: '_receiveBillingInfo',
    UPDATE_COUPON_VALUE: '_updateCouponValue'
  },
  initialize: function() {
    var D = new Date();
    var month = 1;
    var year = D.getFullYear();
    this.billforwardId = '';
    this.accountInfo = { // Billing profile account
      account_code: '',
      username: '',
      email: '',
      first_name: '',
      last_name: '',
      company_name: '',
      hasError: false,
      newBilling: true
    };
    this.billingInfo = { // Billing card information
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
    this.card = {
      number: '',
      cvv: '',
      month: month,
      year: year,
      last_four: null,
      type: '',
      coupon_code: '',
      coupon: 0
    };
    this.errorMessage = '';
    this.fieldErrors = {
      number: false,
      expiry: false,
      cvv: false,
      coupon_code: false,
      first_name: false,
      last_name: false,
      address1: false,
      country: false,
      state: false,
      zip: false,
      city: false,
      month: false,
      year: false
    };
    this.STATUS = STATUS.DEFAULT;
  },
  _accountExists: function() {
    this.accountInfo.newBilling = false;
    this.emitChange();
  },
  _billingInfoExists: function() {
    this.billingInfo.newBilling = false;
    this.emitChange();
  },
  _clearStore: function() {
    this.initialize();
  },
  _receiveBillingInfo: function(payload) {
    this.initialize();
    if (payload.billingInfo && payload.billingInfo.last_four) {
      var cardInfo = {
        last_four: payload.billingInfo.last_four,
        type: payload.billingInfo.card_type,
        month: payload.billingInfo.month,
        year: payload.billingInfo.year
      };
    }
    _.merge(this.billingInfo, payload.billingInfo);
    _.merge(this.accountInfo, payload.accountInfo);
    _.merge(this.card, cardInfo);
    this.billforwardId = payload.billforwardId;
    this.emitChange();
  },
  _submitErrors: function(message) {
    this.STATUS = STATUS.FORM_ERROR;
    this.errorMessage = message;
    this.emitChange();
  },
  _submitSuccess: function() {
    this.STATUS = STATUS.SUCCESS;
    this.errorMessage = '';
    this.emitChange();
  },
  _submitStart: function() {
    this.STATUS = STATUS.ATTEMPTING;
    this.errorMessage = '';
    this.emitChange();
  },
  _updateBillingInfoForm: function({ field, fieldKey, fieldValue }) {
    if (field === 'billing') {
      this.billingInfo[fieldKey] = fieldValue;
    } else if (field === 'account') {
      this.accountInfo[fieldKey] = fieldValue;
    } else if (field === 'card') {
      this.card[fieldKey] = fieldValue;
    }
    this.emitChange();
  },
  _updateErrors: function(hasError) {
    this.STATUS = STATUS.FORM_ERROR;
    _.merge(this.fieldErrors, hasError.fieldErrors);
    _.merge(this.accountInfo, hasError.accountErr);
    this.errorMessage = 'Please make sure all fields are valid.';
    this.emitChange();
  },
  _updateRecurlyErrors: function(error) {
    const errorFields = error.fields;
    debug('Recurly Form errors', errorFields);
    const fieldErrors = {
      number: _.includes(errorFields, 'number'),
      expiry: _.includes(errorFields, 'month') || _.includes(errorFields, 'year'),
      cvv: _.includes(errorFields, 'cvv'),
      first_name: _.includes(errorFields, 'first_name'),
      last_name: _.includes(errorFields, 'last_name')
    };
    _.merge(this.fieldErrors, fieldErrors);
    this.STATUS = STATUS.FORM_ERROR;
    this.errorMessage = error.message;
    this.emitChange();
  },
  _updateCouponValue: function(value) {
    this.card.coupon = value;
    this.emitChange();
  },
  getState: function() {
    return {
      billforwardId: this.billforwardId,
      accountInfo: this.accountInfo,
      billingInfo: this.billingInfo,
      card: this.card,
      errorMessage: this.errorMessage,
      fieldErrors: this.fieldErrors,
      STATUS: this.STATUS
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.billforwardId = state.billforwardId;
    this.accountInfo = state.accountInfo;
    this.billingInfo = state.billingInfo;
    this.card = state.card;
    this.errorMessage = state.errorMessage;
    this.fieldErrors = state.fieldErrors;
    this.STATUS = state.STATUS;
  }
});

module.exports = BillingInfoFormStore;
