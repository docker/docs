'use strict';

import createStore from 'fluxible/addons/createStore';
import { DEFAULT,
         ERROR } from 'stores/enterprisetrialsuccessstore/Constants';
const debug = require('debug')('EnterpriseTrialSuccessStore');

export default createStore({
  storeName: 'EnterpriseTrialSuccessStore',
  handlers: {
    RECEIVE_TRIAL_LICENSE_FACEPALM: '_facepalm',
    RECEIVE_TRIAL_LICENSE: '_receiveTrialLicense',
    RECEIVE_TRIAL_LICENSE_BAD_REQUEST: '_receiveTrialLicenseBadRequest'
  },
  initialize: function() {
    this.license = {};
    this.STATUS = DEFAULT;
  },
  _facepalm: function(err) {
    this.STATUS = ERROR;
    debug(err);
    this.emitChange();
  },
  _receiveTrialLicense: function(license) {
    this.license = license;
    this.emitChange();
  },
  _receiveTrialLicenseBadRequest: function(err) {
    this.STATUS = ERROR;
    debug(err);
    this.emitChange();
  },
  getState: function() {
    return {
      license: this.license,
      STATUS: this.STATUS
    };
  },
  rehydrate: function(state) {
    this.license = state.license;
    this.STATUS = state.STATUS;
  },
  dehydrate: function() {
    return this.getState();
  }
});

