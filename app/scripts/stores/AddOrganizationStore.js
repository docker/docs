'use strict';

var createStore = require('fluxible/addons/createStore');
import { STATUS } from './addorganizationstore/Constants';
var debug = require('debug')('AddOrganizationStore');
var _ = require('lodash');

var noErrorObj = {
  hasError: false,
  error: ''
};

export default createStore({
  storeName: 'AddOrganizationStore',
  handlers: {
    ADD_ORG_CLEAR_FORM: '_clearForm',
    ADD_ORG_UPDATE_FIELD_WITH_VALUE: '_updateFieldWithValue',
    ADD_ORG_ATTEMPT_START: '_addOrgAttemptStart',
    ADD_ORG_BAD_REQUEST: '_badRequest',
    ADD_ORG_SUCCESS: '_addOrgSuccess',
    ADD_ORG_FACEPALM: '_facepalm',
    CREATED_ORGANIZATION: '_clearForm',
    CLEAR_ERRORS: '_clearErrors'
  },
  initialize() {
    this.STATUS = STATUS.DEFAULT;

    this.fields = {
      gravatar_email: {},
      orgname: {},
      location: {},
      company: {},
      profile_url: {}
    };

    this.values = {
      gravatar_email: '',
      orgname: '',
      location: '',
      company: '',
      profile_url: ''
    };
  },
  _facepalm() {
    // this happens if things are screwed and we can't recover gracefully
    this.STATUS = STATUS.FACEPALM;
    this.emitChange();
  },
  _clearForm() {
    this.initialize();
    this.emitChange();
  },
  _updateFieldWithValue({fieldKey, fieldValue}) {
    this._clearErrors(fieldKey);
    this.values[fieldKey] = fieldValue;
    this.emitChange();
  },
  _attemptStart() {
    this.STATUS = STATUS.ATTEMPTING;
    this.emitChange();
  },
  _signupSuccess() {
    this.STATUS = STATUS.SUCCESSFUL_SIGNUP;
    this.emitChange();
  },
  _clearErrors(fieldKey) {
    if (this.STATUS === STATUS.BAD_REQUEST || this.STATUS === STATUS.FACEPALM) {
      this.fields[fieldKey] = noErrorObj;
      this.STATUS = STATUS.DEFAULT;
      this.emitChange();
    }
  },
  _badRequest(obj) {
    /**
     * This function expects keys which match the `this.fields` keys
     * with an array of errors:
     *
     * {
     *   orgname: ['this field is required']
     * }
     */
    let shouldEmitChange = false;
    this.STATUS = STATUS.BAD_REQUEST;

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
