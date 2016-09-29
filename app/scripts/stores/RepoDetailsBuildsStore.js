'use strict';

import createStore from 'fluxible/addons/createStore';
import find from 'lodash/collection/find';
import map from 'lodash/collection/map';
import assign from 'lodash/object/assign';

const debug = require('debug')('RepoDetailsBuildsStore');

const RepoDetailsBuildsStore = createStore({
  storeName: 'RepoDetailsBuildsStore',
  handlers: {
    RECEIVE_BUILD_HISTORY_FOR_REPOSITORY: '_receiveBuilds',
    CANCEL_BUILD_START: '_cancelBuildStart',
    CANCEL_BUILD_SUCCESS: '_cancelBuildSuccess',
    CANCEL_BUILD_ERROR: '_cancelBuildError'
  }
  ,
  initialize() {
    this.results = [];
    this.canceling = {};
    this.count = 0;
  },

  _cancelBuildStart(id) {
    this.canceling = {
      ...this.canceling,
      [id]: 'queued'
    };
    this.emitChange();
  },

  _cancelBuildSuccess(id) {
    this.canceling = {
      ...this.canceling,
      [id]: 'success'
    };
    this.emitChange();
  },

  _cancelBuildError({ id, detail }) {
    this.canceling = {
      ...this.canceling,
      [id]: 'failed'
    };
    this.emitChange();
  },

  _receiveBuilds(res) {
    debug('receiving builds', res);
    this.results = res.results;
    this.count = res.count;
    this.emitChange();
  },

  getState: function() {
    let results = this.results;
    if (this.canceling !== undefined) {
      results = map(this.results, (v) => assign(v, { canceling: this.canceling[v.id] }));
    }
    return {
      count: this.count,
      canceling: this.canceling,
      results
    };
  },

  rehydrate(state) {
    this.results = state.results;
    this.count = state.count;
    this.canceling = state.canceling || {};
  },

  dehydrate() {
    return this.getState();
  }
});

export default RepoDetailsBuildsStore;
