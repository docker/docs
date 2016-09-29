'use strict';

var createStore = require('fluxible/addons/createStore');
import { STATUS } from './common/Constants';
var debug = require('debug')('AccountInfoFormStore');
var _ = require('lodash');

var noErrorObj = {
  hasError: false,
  error: ''
};

export default createStore({
  storeName: 'AccountInfoFormStore',
  handlers: {
    ACCOUNT_INFO_CLEAR_FORM: '_clearForm',
    ACCOUNT_INFO_UPDATE_FIELD_WITH_VALUE: '_updateFieldWithValue',
    ACCOUNT_INFO_ATTEMPT_START: '_attemptStart',
    ACCOUNT_INFO_BAD_REQUEST: '_badRequest',
    ACCOUNT_INFO_SUCCESS: '_success',
    ACCOUNT_INFO_STATUS_CLEAR: '_clearStatus',
    ACCOUNT_INFO_FACEPALM: '_facepalm',
    RECEIVE_USER: '_receiveUser'
  },
  initialize() {
    this.STATUS = STATUS.DEFAULT;

    this.fields = {
      full_name: {},
      company: {},
      location: {},
      profile_url: {},
      gravatar_email: {}
    };

    this.values = {
      full_name: '',
      company: '',
      location: '',
      profile_url: '',
      gravatar_email: ''
    };
  },
  _receiveUser(user) {
    debug('receive user', user);
    this.values.full_name = user.full_name;
    this.values.company = user.company;
    this.values.location = user.location;
    this.values.profile_url = user.profile_url;
    this.values.gravatar_email = user.gravatar_email;
    this.emitChange();
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
    this.STATUS = STATUS.DEFAULT;
    this.values[fieldKey] = fieldValue;
    this.emitChange();
  },
  _attemptStart() {
    this.STATUS = STATUS.ATTEMPTING;
    this.fields = {
      full_name: {},
      company: {},
      location: {},
      profile_url: {},
      gravatar_email: {}
    };
    this.emitChange();
  },
  _success() {
    this.STATUS = STATUS.SUCCESSFUL;
    this.emitChange();
  },
  _clearStatus() {
    this.STATUS = STATUS.DEFAULT;
    this.emitChange();
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
    this.STATUS = STATUS.ERROR;

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
    return this.getState();
  },
  rehydrate(state) {
    debug('rehydrate', state);
    this.fields = state.fields;
    this.values = state.values;
    this.STATUS = state.STATUS;
  }
});
