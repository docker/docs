'use strict';
import createStore from 'fluxible/addons/createStore';
import forEach from 'lodash/collection/forEach';
import isString from 'lodash/lang/isString';
import { STATUS } from './common/Constants';

var AutobuildConfigStore = createStore({
  storeName: 'AutobuildConfigStore',
  handlers: {
    ATTEMPTING_AUTOBUILD_CREATION: '_autobuildCreateAttempt',
    AUTOBUILD_ERROR: '_autobuildConfigError',
    AUTOBUILD_BAD_REQUEST: '_autobuildBadRequest',
    AUTOBUILD_UNAUTHORIZED: '_autobuildUnauthorized',
    AUTOBUILD_SUCCESS: '_autobuildSuccess',
    AUTOBUILD_FORM_UPDATE_FIELD_WITH_VALUE: '_updateFormField',
    SELECT_SOURCE_REPO: '_selectSourceRepo',
    CLEAR_AUTOBUILD_FORM_ERRORS: '_clearErrorStates',
    INITIALIZE_AUTOBUILD_FORM: '_initializeForm',
    RECEIVE_PRIVATE_REPOSTATS: '_getPrivateDefault'
  },
  initialize: function() {
    this.name = '';
    this.namespace = '';
    this.description = '';
    this.isPrivate = 'public';
    this.provider = '';
    this.sourceRepoName = '';
    this.active = true;
    this.error = {};
    this.success = '';
    this.STATUS = STATUS.DEFAULT;
  },
  _autobuildConfigError: function(err) {
    //TODO: handle config error here
    this.error.general = 'An error occurred while configuring your automated build. Please try again later.';
    setTimeout(this._clearErrorStates.bind(this), 5000);
    this.emitChange();
  },
  _autobuildCreateAttempt: function() {
    this.error.buildTags = '';
    this.STATUS = STATUS.ATTEMPTING;
    this.emitChange();
  },
  _autobuildBadRequest: function(err) {
    forEach(err, (val, key) => {
      this.error[key] = val.toString();
    });

    //For build_tags, make it a global error
    if (err.build_tags) {
      this.error.buildTags = 'Invalid character(s) provided in build tags configuration. Please check your input.';
    }

    if (err.detail || isString(err)) {
      this.error.detail = err.detail || err;
    }

    this.STATUS = STATUS.DEFAULT;
    this.emitChange();
  },
  _autobuildUnauthorized: function(err) {
    this.error.general = 'You have no permissions to create an automated build in this namespace.';
    setTimeout(this._clearErrorStates.bind(this), 5000);
    this.STATUS = STATUS.DEFAULT;
    this.emitChange();
  },
  _autobuildSuccess: function(err) {
    this.success = 'Successfully configured an automated build repository.';
    this.STATUS = STATUS.SUCCESSFUL;
    setTimeout(this._clearErrorStates.bind(this), 5000);
    this.emitChange();
  },
  _clearErrorStates: function() {
    this.error = {};
    this.success = '';
    this.STATUS = STATUS.DEFAULT;
    this.emitChange();
  },
  _getPrivateDefault: function(stats) {
    this.isPrivate = stats.default_repo_visibility;
  },
  _initializeForm: function({ name, namespace }) {
    this.name = name;
    this.namespace = namespace;
    this.description = '';
  },
  _selectSourceRepo: function(repo) {
    this.sourceRepoName = repo.full_name;
    this.emitChange();
  },
  _updateFormField: function({fieldKey, fieldValue}) {
    this[fieldKey] = fieldValue;
    if (fieldKey === 'name' || fieldKey === 'namespace') {
      delete this.error.dockerhub_repo_name;
    }
    if (fieldKey === 'description') {
      delete this.error.description;
    }
    delete this.error.detail;
    this.STATUS = STATUS.DEFAULT;
    this.emitChange();
  },
  getState: function() {
    return {
      name: this.name,
      namespace: this.namespace,
      description: this.description,
      isPrivate: this.isPrivate,
      provider: this.provider,
      sourceRepoName: this.sourceRepoName,
      active: this.active,
      error: this.error,
      success: this.success,
      STATUS: this.STATUS
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.name = state.name;
    this.namespace = state.namespace;
    this.description = state.description;
    this.isPrivate = state.isPrivate;
    this.provider = state.provider;
    this.sourceRepoName = state.sourceRepoName;
    this.active = state.active;
    this.error = state.error;
    this.success = state.success;
    this.STATUS = state.STATUS;
  }
});

module.exports = AutobuildConfigStore;
