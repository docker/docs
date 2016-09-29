'use strict';

import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('RepoDetailsBuildLogsStore');

export default createStore({
  storeName: 'RepoDetailsBuildLogsStore',
  handlers: {
    BUILD_LOGS_RECEIVE: '_receiveBuildLogs'
  },
  initialize() {
    this.build_results = {};
  },
  _receiveBuildLogs(res) {
    this.build_results = res.build_results;
    this.emitChange();
  },
  getState() {
    return {
      build_results: this.build_results
    };
  },
  rehydrate(state) {
    this.build_results = state.build_results;
  },
  dehydrate() {
    return this.getState();
  }
});
