'use strict';

import createStore from 'fluxible/addons/createStore';
const debug = require('debug')('PipelineHistory');

export default createStore({
  storeName: 'PipelineHistoryStore',
  handlers: {
    RECEIVE_PIPELINE_HISTORY: '_receivePipelineHistory'
  },
  initialize() {
    this.results = {};
  },
  _receivePipelineHistory(data) {
    this.results[data.slug] = {};
    this.results[data.slug].results = data.payload.results;
    this.results[data.slug].count = data.payload.count;
    this.emitChange();
  },
  getState() {
    return {
      results: this.results
    };
  },
  dehydrate() {
    return this.getState();
  },
  rehydrate(state) {
    this.results = state.results;
  }
});
