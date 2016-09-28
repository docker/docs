'use strict';

import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('RepoDetailsBuildsStore');

var RepoDetailsBuildsStore = createStore({
  storeName: 'RepoDetailsBuildsStore',
  handlers: {
    RECEIVE_BUILD_HISTORY_FOR_REPOSITORY: '_receiveBuilds'
  },
  initialize() {
    this.results = [];
    this.count = 0;
  },
  _receiveBuilds(res) {
    debug('receiving builds', res);
    this.results = res.results;
    this.count = res.count;
    this.emitChange();
  },
  getState: function() {
      return {
        results: this.results,
        count: this.count
      };
  },
  rehydrate(state) {
    this.results = state.results;
    this.count = state.count;
  },
  dehydrate() {
    return this.getState();
  }
});

module.exports = RepoDetailsBuildsStore;
