'use strict';

import _ from 'lodash';
import { STATUS } from './loginstore/Constants';
import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('LoginStore');

export default createStore({
  storeName: 'LoginStore',
  handlers: {
    LOGIN_ATTEMPT_START: '_loginAttemptStart',
    LOGIN_UNAUTHORIZED: '_loginUnauthorized',
    LOGIN_UNAUTHORIZED_DETAIL: '_loginUnauthorizedDetail',
    LOGIN_BAD_REQUEST: '_badRequest',
    LOGIN_ERROR: '_loginError',
    LOGIN_UPDATE_FIELD_WITH_VALUE: '_updateFieldWithValue',
    LOGIN_CLEAR: '_clearLoginForm'
  },
  initialize() {
    this.STATUS = STATUS.DEFAULT;
    this.globalFormError = '';

    this.fields = {
      username: {},
      password: {}
    };

    this.values = {
      username: '',
      password: ''
    };
  },
  _clearLoginForm() {
    debug('Clearing');
    this.initialize();
    this.emitChange();
  },
  _loginAttemptStart() {
    this.STATUS = STATUS.ATTEMPTING_LOGIN;
    this.emitChange();
  },
  _loginError(err){
    this.STATUS = STATUS.GENERIC_ERROR;
    this.globalFormError = 'There was an error contacting the server. Please try again later.';
    this.emitChange();
  },
  _badRequest(obj) {
    this.STATUS = STATUS.DEFAULT;
    /**
     * This function expects keys which match the `this.fields` keys
     * with an array of errors:
     *
     * {
     *   username: ['this field is required']
     * }
     */
    let shouldEmitChange = false;

    // cycle through the possible form fields
    this.fields = _.mapValues(this.fields, function (errorObject, key) {
      if(_.has(obj, key)) {
        shouldEmitChange = true;
        return {
          hasError: !!obj[key],
          error: obj[key][0]
        };
      } else {
        return errorObject;
      }
    });

    if(shouldEmitChange) {
      this.emitChange();
    }
  },
  _loginUnauthorized() {
    this.STATUS = STATUS.ERROR_UNAUTHORIZED;
    this.globalFormError = 'Login Failed. The username or password may be incorrect.';
    this.emitChange();
  },
  _loginUnauthorizedDetail({detail}) {
    this.STATUS = STATUS.ERROR_UNAUTHORIZED;
    this.globalFormError = detail;
    this.emitChange();
  },
  getState() {
    return {
      fields: this.fields,
      values: this.values,
      STATUS: this.STATUS,
      globalFormError: this.globalFormError
    };
  },
  _updateFieldWithValue: function({fieldKey, fieldValue}){
    this.fields[fieldKey] = {
      hasError: false,
      error: ''
    };
    this.values[fieldKey] = fieldValue;
    this.emitChange();
  },
  dehydrate: function() {
    return {};
  },
  rehydrate: function(state) {
    this.state = state;
  }
});
