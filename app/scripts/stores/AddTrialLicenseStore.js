'use strict';

import createStore from 'fluxible/addons/createStore';
import { ATTEMPTING_DOWNLOAD,
         BAD_REQUEST,
         DEFAULT,
         FACEPALM,
         SUCCESSFUL_DOWNLOAD } from 'stores/addtriallicensestore/Constants';
const debug = require('debug')('AddTrialLicenseStore');

export default createStore({
  storeName: 'AddTrialLicenseStore',
  handlers: {
    ATTEMPTING_LICENSE_DOWNLOAD_START: '_attemptingLicenseDownloadStart',
    DOWNLOAD_LICENSE_CONTENT_BAD_REQUEST: '_downloadLicenseContentBadRequest',
    DOWNLOAD_LICENSE_CONTENT_FACEPALM: '_facepalm',
    RECEIVE_LICENSE_DOWNLOAD_CONTENT: '_receiveLicenseDownloadContent'
  },
  initialize: function() {
    this.error = '';
    this.STATUS = DEFAULT;
  },
  _attemptingLicenseDownloadStart: function() {
    this.STATUS = ATTEMPTING_DOWNLOAD;
    this.error = '';
    this.emitChange();
  },
  _clearFeedbackStates: function() {
    this.STATUS = DEFAULT;
    this.error = '';
    this.emitChange();
  },
  _downloadLicenseContentBadRequest: function(err) {
    this.STATUS = BAD_REQUEST;
    this.error = err;
    this.emitChange();
  },
  _facepalm: function(err) {
    this.STATUS = FACEPALM;
    debug(err);
    this.error = 'Sorry, an error occured and your license is unavailable at this time.';
    this.emitChange();
  },
  _receiveLicenseDownloadContent: function() {
    this.STATUS = SUCCESSFUL_DOWNLOAD;
    this.error = '';
    setTimeout(this._clearFeedbackStates.bind(this), 5000);
    this.emitChange();
  },
  getState: function() {
    return {
      error: this.error,
      STATUS: this.STATUS
    };
  },
  rehydrate: function(state) {
    this.error = state.error;
    this.STATUS = state.STATUS;
  },
  dehydrate: function() {
    return this.getState();
  }
});

