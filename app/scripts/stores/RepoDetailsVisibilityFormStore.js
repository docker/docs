'use strict';

import createStore from 'fluxible/addons/createStore';
import { STATUS } from './repovisibilitystore/Constants.js';
const debug = require('debug')('RepoDetailsVisibilityFormStore');

var RepoDetailsVisibilityFormStore = createStore({
  storeName: 'RepoDetailsVisibilityFormStore',
  handlers: {
    VISIBILITY_BAD_REQUEST: '_badRequest',
    VISIBILITY_ERROR: '_visibilityError',
    TOGGLE_VISIBILITY_ATTEMPT_START: '_toggleVisibilityAttemptStart',
    TOGGLE_VISIBILITY_SUCCESS: '_toggleSuccess',
    RECEIVE_PRIVATE_REPOSTATS: '_receivePrivateRepoStats',
    RECEIVE_REPOSITORY: '_receiveRepository',
    REPO_DETAILS_VISIBILITY_UPDATE_FIELD_WITH_VALUE: '_updateFieldWithValue',
    TOGGLE_VISIBILITY_REPO_NAME_CONFIRM_BOX: '_toggleConfirmBox'
  },
  initialize: function() {
    this.badRequest = '';
    this.error = '';
    this.success = '';
    this.isPrivate = false;
    this.privateRepoLimit = null;
    this.numPrivateReposAvailable = null;
    this.STATUS = STATUS.DEFAULT;
    this.values = {
      confirmRepoName: ''
    };
  },
  _badRequest: function(res) {
    this.initialize();
    this.badRequest = res.detail;
    this.STATUS = STATUS.FORM_ERROR;
    this.emitChange();
  },
  _clearErrors: function () {
    this.error = '';
    this.badRequest = '';
  },
  _toggleVisibilityAttemptStart: function() {
    this.STATUS = STATUS.ATTEMPTING;
    this.emitChange();
  },
  _visibilityError: function(maybeError) {
    this.STATUS = STATUS.FORM_ERROR;
    if (maybeError) {
      this.error = maybeError.detail;
    } else {
      this.error = 'No private repositories available';
    }
    this.emitChange();
  },
  _toggleSuccess: function(isPrivate) {
    this.initialize();
    this.isPrivate = isPrivate;
    this.emitChange();
  },
  _toggleConfirmBox: function() {
    if (this.STATUS === STATUS.DEFAULT) {
      this.STATUS = STATUS.SHOWING_CONFIRM_BOX;
    } else {
      this.STATUS = STATUS.DEFAULT;
    }
    this._clearErrors();
    this.values.confirmRepoName = '';
    this.emitChange();
  },
  _receivePrivateRepoStats: function(stats) {
    /*eslint-disable camelcase */
    this.numPrivateReposAvailable = stats.private_repo_available;
    this.privateRepoLimit = stats.private_repo_limit;
    /*eslint-enable camelcase */
    this.emitChange();
  },
  _receiveRepository: function(res) {
    this.initialize();
    this.isPrivate = res.is_private;
    this.emitChange();
  },
  _updateFieldWithValue: function({fieldKey, fieldValue}){
    this.values[fieldKey] = fieldValue;
    this._clearErrors();
    this.emitChange();
  },
  getState: function() {
    return {
      badRequest: this.badRequest,
      error: this.error,
      success: this.success,
      isPrivate: this.isPrivate,
      privateRepoLimit: this.privateRepoLimit,
      numPrivateReposAvailable: this.numPrivateReposAvailable,
      values: this.values,
      STATUS: this.STATUS
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.badRequest = state.badRequest;
    this.error = state.error;
    this.success = state.success;
    this.isPrivate = state.isPrivate;
    this.numPrivateReposAvailable = state.numPrivateReposAvailable;
    this.privateRepoLimit = state.privateRepoLimit;
    this.values = state.values;
    this.STATUS = state.STATUS;
  }
});

module.exports = RepoDetailsVisibilityFormStore;
