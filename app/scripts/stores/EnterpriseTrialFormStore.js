'use strict';

import createStore from 'fluxible/addons/createStore';
import { ATTEMPTING, DEFAULT, FACEPALM, SUCCESSFUL_SIGNUP } from 'stores/enterprisetrialstore/Constants';
var debug = require('debug')('EnterpriseTrialFormStore');
var _ = require('lodash');

var noErrorObj = {
  hasError: false,
  error: ''
};

export default createStore({
  storeName: 'EnterpriseTrialFormStore',
  handlers: {
    ENTERPRISE_TRIAL_RECEIVE_ORGS: '_receiveOrgs',
    ENTERPRISE_TRIAL_CLEAR_FORM: '_clearForm',
    ENTERPRISE_TRIAL_UPDATE_FIELD_WITH_VALUE: '_updateFieldWithValue',
    ENTERPRISE_TRIAL_ATTEMPT_START: '_attemptStart',
    ENTERPRISE_TRIAL_BAD_REQUEST: '_badRequest',
    ENTERPRISE_TRIAL_SUCCESS: '_signupSuccess',
    ENTERPRISE_TRIAL_FACEPALM: '_facepalm',
    CREATED_ORGANIZATION: '_clearForm'
  },
  initialize() {
    this.STATUS = DEFAULT;
    this.orgs = [];
    this.globalFormError = '';

    this.fields = {
      firstName: {},
      lastName: {},
      companyName: {},
      jobFunction: {},
      email: {},
      phoneNumber: {},
      country: {},
      state: {},
      namespace: {}
    };

    this.values = {
      namespace: '',
      firstName: '',
      lastName: '',
      jobFunction: '',
      companyName: '',
      email: '',
      phoneNumber: '',
      country: 'US',
      state: ''
    };
  },
  _facepalm() {
    // this happens if things are screwed and we can't recover gracefully
    this.STATUS = FACEPALM;
    this.globalFormError = 'Something went wrong on the server. We have been alerted to this issue';
    this.emitChange();
  },
  _clearForm() {
    this.initialize();
    this.emitChange();
  },
  _clearErrors() {
    this.fields = {
      firstName: {},
      lastName: {},
      jobFunction: {},
      companyName: {},
      email: {},
      phoneNumber: {},
      country: {},
      state: {},
      namespace: {}
    };
    this.globalFormError = '';
  },
  _updateFieldWithValue({fieldKey, fieldValue}) {
    this.STATUS = DEFAULT;
    this.globalFormError = '';
    this.values[fieldKey] = fieldValue;
    this.emitChange();
  },
  _attemptStart() {
    this.STATUS = ATTEMPTING;
    this.emitChange();
  },
  _signupSuccess() {
    this.STATUS = SUCCESSFUL_SIGNUP;
    this._clearErrors();
    this.emitChange();
  },
  _receiveOrgs(namespaces) {
    this.orgs = namespaces;
    //will always have at least current logged in namespace
    this.values.namespace = namespaces[0];
    this.emitChange();
  },
  _badRequest(obj) {
    this._clearErrors();
    this.STATUS = DEFAULT;

    // cycle through the possible form fields
    this.fields = _.mapValues(this.fields, (errorObject, key) => {
      if(_.has(obj, key)) {
        return {
          hasError: !!obj[key],
          error: obj[key][0]
        };
      } else {
        return errorObject;
      }
    });

    if(obj && obj.non_field_errors) {
      this.globalFormError = obj.non_field_errors[0];
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
    this.fields = state.fields;
    this.values = state.values;
    this.STATUS = state.STATUS;
    this.globalFormError = state.globalFormError;

  }
});
