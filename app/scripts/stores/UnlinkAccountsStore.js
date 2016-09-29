'use strict';

import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('UnlinkAccountsStore');

var BitbucketLinkStore = createStore({
  storeName: 'UnlinkAccountsStore',
  handlers: {
    GITHUB_UNLINK_ERROR: _unlinkGithubError,
    BITBUCKET_UNLINK_ERROR: _unlinkBitbucketError
  },
  initialize: function() {
    this.error = '';
  },
  _unlinkGithubError: function() {
    this.error = 'Error unlinking Github Account. Please try again later';
    this.emitChange();
  },
  _unlinkBitbucketError: function() {
    this.error = 'Error unlinking Bitbucket Account. Please try again later';
    this.emitChange();
  },
  getState: function() {
    return {
      error: this.error
    };
  },
  rehydrate: function(state) {
    this.error = state.error;
  },
  dehydrate: function() {
    return this.getState();
  }
});

module.exports = UnlinkAccountsStore;
