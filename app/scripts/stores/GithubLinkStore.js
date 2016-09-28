'use strict';

import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('GithubLinkStore');

var GithubLinkStore = createStore({
  storeName: 'GithubLinkStore',
  handlers: {
    RECEIVE_GITHUB_ID: '_receiveID',
    GITHUB_ID_ERROR: '_idError',
    GITHUB_SECURITY_ERROR: '_githubSecurityError',
    GITHUB_ASSOCIATE_ERROR: '_githubAssociateError'
  },
  initialize: function() {
    this.githubClientID = '';
    this.error = '';
  },
  _receiveID: function(res) {
    this.githubClientID = res.client_id;
    this.emitChange();
  },
  _idError: function(err) {
    debug(err);
  },
  _githubAssociateError: function(errorState) {
    debug(errorState);
    if (errorState.detail) {
        this.error = errorState.detail;
    } else {
        this.error = 'There was an error during the Github account link. Please check that you do not have the same Github account linked to another Docker Hub account.';
    }
    this.emitChange();
    setTimeout(this._clearError.bind(this), 5000);
  },
  _githubSecurityError: function(errorState) {
    debug(errorState);
    this.error = 'There was a security error during the github account linking process.';
    this.emitChange();
    setTimeout(this._clearError.bind(this), 5000);
  },
  _clearError: function() {
    this.error = '';
    this.emitChange();
  },
  getState: function() {
    return {
      githubClientID: this.githubClientID,
      error: this.error
    };
  },
  rehydrate: function(state) {
    this.githubClientID = state.githubClientID;
    this.error = state.error;
  },
  dehydrate: function() {
    return this.getState();
  }
});

module.exports = GithubLinkStore;
