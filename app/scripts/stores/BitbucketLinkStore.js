'use strict';

import _ from 'lodash';
import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('BitbucketLinkStore');

var BitbucketLinkStore = createStore({
  storeName: 'BitbucketLinkStore',
  handlers: {
    RECEIVE_BITBUCKET_AUTH_URL: '_receiveUrl',
    BITBUCKET_AUTH_URL_ERROR: '_urlError',
    BITBUCKET_ASSOCIATE_ERROR: '_associateError'
  },
  initialize: function() {
    this.authURL = '';
    this.error = '';
  },
  _associateError: function(body) {
    debug(body);
    if (_.has(body, 'detail') && _.isString(body.detail)) {
        this.error = body.detail;
    } else {
        this.error = 'Error linking your account to Bitbucket. Please check that you do not have the same Bitbucket account linked to another Docker Hub account.';
    }
    this.emitChange();
    setTimeout(this._clearError.bind(this), 5000);
  },
  _receiveUrl: function(res) {
    this.authURL = res.bitbucket_authorization_url;
    this.emitChange();
  },
  _urlError: function(err) {
    debug(err);
    this.error = 'Error linking your account to bitbucket.';
    this.emitChange();
    setTimeout(this._clearError.bind(this), 5000);
  },
  _clearError: function() {
    this.error = '';
    this.emitChange();
  },
  setURL: function(url) {
    this.authURL = url;
  },
  getState: function() {
    return {
      authURL: this.authURL,
      error: this.error
    };
  },
  rehydrate: function(state) {
    this.authURL = state.authURL;
    this.error = state.error;
  },
  dehydrate: function() {
    return this.getState();
  }
});

module.exports = BitbucketLinkStore;
