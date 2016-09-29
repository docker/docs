'use strict';

import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('PrivateReposUsageStore');

var PrivateRepoUsageStore = createStore({
  storeName: 'PrivateRepoUsageStore',
  handlers: {
    RECEIVE_PRIVATE_REPOSTATS: '_receivePrivateRepoStats',
    PRIVATE_REPOSTATS_NO_PERMISSIONS: '_notAvailable'
  },
  initialize: function() {
    this.privateRepoUsed = 0;
    this.numFreePrivateRepos = 0;
    this.defaultRepoVisibility = 'public';
    this.privateRepoAvailable = 0;
    this.privateRepoPercentUsed = 0;
    this.privateRepoLimit = 0;
    this.notAvailable = false;
  },
  _receivePrivateRepoStats: function(stats) {
    this.notAvailable = false;
    /*eslint-disable camelcase */
    this.privateRepoUsed = stats.private_repo_used;
    this.numFreePrivateRepos = stats.num_free_private_repos;
    this.defaultRepoVisibility = stats.default_repo_visibility;
    this.privateRepoAvailable = stats.private_repo_available;
    this.privateRepoPercentUsed = stats.private_repo_percent_used;
    this.privateRepoLimit = stats.private_repo_limit;
    /*eslint-enable camelcase */
    this.emitChange();
  },
  _notAvailable: function(err) {
    //No permissions to see the private repo stats for this org
    this.notAvailable = true;
    this.emitChange();
  },
  getState: function() {
    return {
      privateRepoUsed: this.privateRepoUsed,
      numFreePrivateRepos: this.numFreePrivateRepos,
      defaultRepoVisibility: this.defaultRepoVisibility,
      privateRepoAvailable: this.privateRepoAvailable,
      privateRepoPercentUsed: this.privateRepoPercentUsed,
      privateRepoLimit: this.privateRepoLimit,
      notAvailable: this.notAvailable
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.privateRepoUsed = state.privateRepoUsed;
    this.numFreePrivateRepos = state.numFreePrivateRepos;
    this.defaultRepoVisibility = state.defaultRepoVisibility;
    this.privateRepoAvailable = state.privateRepoAvailable;
    this.privateRepoPercentUsed = state.privateRepoPercentUsed;
    this.privateRepoLimit = state.privateRepoLimit;
    this.notAvailable = state.notAvailable;
  }
});

module.exports = PrivateRepoUsageStore;
