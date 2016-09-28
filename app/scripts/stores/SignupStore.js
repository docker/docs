'use strict';

var createStore = require('fluxible/addons/createStore');
import { STATUS } from './signupstore/Constants';
var debug = require('debug')('SignupStore');
var _ = require('lodash');

var noErrorObj = {
  hasError: false,
  error: ''
};

export default createStore({
  storeName: 'SignupStore',
  handlers: {
    SIGNUP_CLEAR_FORM: '_signupClearForm',
    SIGNUP_CLEAR_PASSWORD: '_signupClearPassword',
    SIGNUP_UPDATE_FIELD_WITH_VALUE: '_updateFieldWithValue',
    SIGNUP_ATTEMPT_START: '_signupAttemptStart',
    SIGNUP_BAD_REQUEST: '_badRequest',
    SIGNUP_SUCCESS: '_signupSuccess'
  },
  initialize() {
    this.STATUS = STATUS.DEFAULT;

    this.fields = {
      username: {},
      email: {},
      password: {}
    };

    this.values = {
      username: '',
      email: '',
      password: ''
    };
  },
  _signupClearForm() {
    this.initialize();
    this.emitChange();
  },
  _signupClearPassword() {
    this.values.password = '';
    this.emitChange();
  },
  _updateFieldWithValue({fieldKey, fieldValue}){
    this.values[fieldKey] = fieldValue;
    if (fieldValue) {
      this.fields[fieldKey] = this._validate({fieldKey, fieldValue});
    }
    this.emitChange();
  },
  _signupAttemptStart() {
    debug('attempting Signup');
    this.STATUS = STATUS.ATTEMPTING_SIGNUP;
    this.emitChange();
  },
  _signupSuccess() {
    this.STATUS = STATUS.SUCCESSFUL_SIGNUP;
    this.emitChange();
  },
  _badRequest(obj) {
    let shouldEmitChange = false;

    // cycle through the possible form fields
    _.forEach(_.keys(this.fields),
              (key) => {
                if(_.has(obj, key)) {
                  shouldEmitChange = true;
                  var newField = {};
                  newField.hasError = !!obj[key];
                  newField.error = obj[key][0];
                  this.fields[key] = newField;
                }
              });
    if(shouldEmitChange) {
      this.emitChange();
    }
  },
  validations: {
    username(value) {
      if (value.length < 4){
        return {
          hasError: true,
          error: 'Username must be at least four characters long'
        };
      } else if (!/^[A-Za-z0-9]+$/.test(value)) {
        return {
          hasError: true,
          error: 'Username must contain only letters and digits'
        };
      } else {
        return noErrorObj;
      }
    }
  },
  _validate({fieldKey, fieldValue}) {
    if(_.isFunction(this.validations[fieldKey])) {
      return this.validations[fieldKey](fieldValue);
    } else {
      return noErrorObj;
    }
  },
  getState() {
    return {
      fields: this.fields,
      values: this.values,
      STATUS: this.STATUS
    };
  },
  dehydrate() {
    return {};
  },
  rehydrate(state) {
    this.state = state;
  }
});
