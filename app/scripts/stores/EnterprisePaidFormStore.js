'use strict';

import createStore from 'fluxible/addons/createStore';
import { STATUS } from './common/Constants';
var debug = require('debug')('EnterprisePaidFormStore');
import each from 'lodash/collection/each';
import includes from 'lodash/collection/includes';
import merge from 'lodash/object/merge';
import keys from 'lodash/object/keys';
import has from 'lodash/object/has';
import mapValues from 'lodash/object/mapValues';
import isString from 'lodash/lang/isString';

var noErrorObj = {
  hasError: false,
  error: ''
};

export default createStore({
  storeName: 'EnterprisePaidFormStore',
  handlers: {
    ENTERPRISE_PAID_RECEIVE_ORGS: '_receiveOrgs',
    ENTERPRISE_PAID_CLEAR_FORM: '_clearStore',
    ENTERPRISE_PAID_UPDATE_FIELD_WITH_VALUE: '_updateFieldWithValue',

    ENTERPRISE_PAID_ATTEMPT_START: '_enterprisePaidAttemptStart',
    ENTERPRISE_PAID_BAD_REQUEST: '_badRequest',
    ENTERPRISE_PAID_SUCCESS: '_signupSuccess',
    ENTERPRISE_PAID_FACEPALM: '_facepalm',

    BILLING_SUBMIT_START: '_enterprisePaidAttemptStart',
    BILLING_SUBMIT_SUCCESS: '_signupSuccess',
    BILLING_SUBMIT_ERROR: '_badRequest',

    GET_RECURLY_ERROR: '_updateRecurlyErrors',
    ENTERPRISE_PAID_GET_RECURLY_ERROR: 'recurlyError',
    ENTERPRISE_PAID_API_ERROR: '_apiError',
    ENTERPRISE_PAID_ERRORS: '_validateErrors',
    ENTERPRISE_PAID_POPULATE_FORM: '_populateForm'
  },
  initialize() {
    this.STATUS = STATUS.DEFAULT;

    this.globalFormError = '';
    this.orgs = [];

    this.fields = {
      first_name: {},
      last_name: {},
      postal_code: {},
      number: {},
      month: {},
      year: {},
      cvv: {},
      address1: {},
      city: {},
      state: {},
      country: {},
      expiry: {},
      email: {}
    };

    this.values = {
      first_name: '',
      last_name: '',
      postal_code: '',
      number: '',
      month: '01',
      year: '2015',
      cvv: '',
      address1: '',
      city: '',
      state: '',
      country: 'US',
      last_four: '',
      card_type: '',
      account_first: '',
      account_last: '',
      company_name: '',
      email: ''
    };
  },
  _clearStore(){
    this.initialize();
    this.emitChange();
  },
  _populateForm({
    first_name,
    last_name,
    zip,
    month,
    year,
    address1,
    address2,
    city,
    state,
    country,
    last_four,
    card_type,
    account_first,
    account_last,
    company_name,
    email
    }) {
    var D = new Date();
    var defaultMonth = D.getMonth();
    var defaultYear = D.getFullYear() + 1;
    var defaultCountry = 'US';
    this.fields = {
      first_name: {},
      last_name: {},
      postal_code: {},
      number: {},
      month: {},
      year: {},
      cvv: {},
      address1: {},
      address2: {},
      city: {},
      state: {},
      country: {},
      expiry: {},
      email: {}
    };
    this.values = {
      first_name,
      last_name,
      postal_code: zip,
      month: month || defaultMonth,
      year: year || defaultYear,
      address1,
      address2,
      city,
      state,
      country: country || defaultCountry,
      last_four,
      card_type,
      account_first,
      account_last,
      company_name,
      email
    };
    this.STATUS = STATUS.DEFAULT;
    this.globalFormError = '';
    this.emitChange();
  },
  _updateRecurlyErrors(error) {
    const errorFields = error.fields;
    debug('Recurly Form errors', errorFields);
    let fieldErrors = {
      number: {
        hasError: includes(errorFields, 'number'),
        error: 'There was an error processing your card'
      },
      expiry: {
        hasError: includes(errorFields, 'month') || includes(errorFields, 'year'),
        error: 'This field is invalid'
      },
      cvv: {
        hasError: includes(errorFields, 'cvv'),
        error: 'This field is invalid'
      },
      first_name: {
        hasError: includes(errorFields, 'first_name'),
        error: 'This field is required'
      },
      last_name: {
        hasError: includes(errorFields, 'last_name'),
        error: 'This field is required'
      },
      postal_code: {
        hasError: includes(errorFields, 'postal_code'),
        error: 'This field is invalid'
      }
    };
    merge(this.fields, fieldErrors);
    this.STATUS = STATUS.DEFAULT;
    this.globalFormError = error.message;
    this.emitChange();
  },
  _facepalm() {
    // this happens if things are screwed and we can't recover gracefully
    this.STATUS = STATUS.FACEPALM;
    this.emitChange();
  },
  recurlyError(fields) {
    this.STATUS = STATUS.DEFAULT;
    var emitChange = false;
    each(fields, function(val, idx) {
      if(includes(keys(this.fields), val)) {
        emitChange = true;
        this.fields[val] = {
          hasError: true,
          error: 'This field is required'
        };
      }
    }, this);

    if(emitChange) {
      this.emitChange();
    }
  },
  _validateErrors(hasError) {
    each(hasError, (v, k) => {
      if (v) {
        if (includes(['number', 'cvv', 'month', 'year', 'country'], k)) {
          this.fields[k] = {
            hasError: true,
            error: 'Invalid ' + k
          };
        } else {
          this.fields[k] = {
            hasError: true,
            error: 'Required'
          };
        }
      }
    });
    this.emitChange();
  },
  _updateFieldWithValue({fieldKey, fieldValue}) {
    debug(fieldKey, fieldValue);
    this.fields[fieldKey] = {hasError: false, error: ''};
    if (includes(['month', 'year'], fieldKey)) {
      this.fields.expiry = {hasError: false, error: ''};
    }
    if (fieldKey === 'number') {
      let card_type = window.recurly.validate.cardType(fieldValue);
      this.values.card_type = card_type;
    }
    this.values[fieldKey] = fieldValue;
    this.emitChange();
  },
  _signupSuccess() {
    this.STATUS = STATUS.SUCCESSFUL;
    this.emitChange();
  },
  _receiveOrgs(namespaces) {
    this.orgs = namespaces;
    this.emitChange();
  },
  _apiError() {
    this.STATUS = STATUS.FACEPALM;
    this.emitChange();
  },
  _enterprisePaidAttemptStart() {
    this.STATUS = STATUS.ATTEMPTING;
    this.emitChange();
  },
  _badRequest(obj) {
    this.STATUS = STATUS.ERROR;

    // cycle through the possible form fields
    this.fields = mapValues(this.fields, function (errorObject, key) {
      if(has(obj, key)) {
        return {
          hasError: !!obj[key],
          error: obj[key][0]
        };
      } else {
        return errorObject;
      }
    });

    if(has(obj, 'non_field_errors')) {
      this.globalFormError = obj.non_field_errors[0];
    } else if (has(obj, 'detail')) {
      this.globalFormError = obj.detail;
    } else if (isString(obj)) {
      this.globalFormError = obj;
    }

    this.emitChange();
  },
  getState() {
    return {
      fields: this.fields,
      values: this.values,
      STATUS: this.STATUS,
      orgs: this.orgs,
      globalFormError: this.globalFormError
    };
  },
  dehydrate() {
    return this.getState();
  },
  rehydrate(state) {
    this.orgs = state.orgs;
    this.values = state.values;
    this.fields = state.fields;
    this.STATUS = state.STATUS;
    this.globalFormError = this.globalFormError;
  }
});
