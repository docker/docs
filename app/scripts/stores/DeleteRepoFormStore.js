'use strict';

import createStore from 'fluxible/addons/createStore';
import { STATUS } from './deleterepostore/Constants.js';
const debug = require('debug')('DeleteRepoFormStore');

var DeleteRepoFormStore = createStore({
  storeName: 'DeleteRepoFormStore',
  handlers: {
    DELETE_REPO_ATTEMPT_START: '_deleteRepoAttemptStart',
    DELETE_REPO_BAD_REQUEST: '_deleteRepoBadRequest',
    DELETE_REPO_ERROR: '_deleteRepoError',
    DELETE_REPO_UPDATE_FIELD_WITH_VALUE: '_updateFieldWithValue',
    RECEIVE_REPOSITORY: '_receiveRepository',
    TOGGLE_DELETE_REPO_NAME_CONFIRM_BOX: '_toggleConfirmBox'
  },
  initialize: function() {
    this.error = '';
    this.STATUS = STATUS.DEFAULT;
    this.values = {
      confirmRepoName: ''
    };
  },
  _deleteRepoAttemptStart: function() {
    this.STATUS = STATUS.ATTEMPTING;
    this.emitChange();
  },
  _deleteRepoBadRequest: function(res) {
    this.STATUS = STATUS.FORM_ERROR;
    this.error = res.detail ? res.detail
                            : 'Error deleting repository. Please verify if you have permissions.';
    this.emitChange();
  },
  _deleteRepoError: function() {
    this.STATUS = STATUS.FORM_ERROR;
    this.error = 'Error deleting repository. Please verify if you have permissions.';
    this.emitChange();
  },
   _receiveRepository: function() {
    this.initialize();
    this.emitChange();
  },
  _toggleConfirmBox: function() {
    if (this.STATUS === STATUS.DEFAULT) {
      this.STATUS = STATUS.SHOWING_CONFIRM_BOX;
    } else {
      this.STATUS = STATUS.DEFAULT;
    }
    this.error = '';
    this.values.confirmRepoName = '';
    this.emitChange();
  },
  _updateFieldWithValue: function({fieldKey, fieldValue}){
    this.values[fieldKey] = fieldValue;
    this.error = '';
    this.emitChange();
  },
  getState: function() {
    return {
      error: this.error,
      STATUS: this.STATUS,
      values: this.values
    };
  },
  rehydrate: function(state) {
    this.error = state.error;
    this.STATUS = state.STATUS;
    this.values = state.values;
  },
  dehydrate: function() {
    return this.getState();
  }
});

module.exports = DeleteRepoFormStore;
