'use strict';

import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('DashboardContribsStore');

export default createStore({
  storeName: 'DashboardContribsStore',
  handlers: {
    RECEIVE_CONTRIB: '_receiveContribRepos',
    LOGOUT: 'initialize'
  },
  initialize: function() {
    this.count = 0;
    this.contribs = [];
    this.next = null;
    this.prev = null;
  },
  _receiveContribRepos: function(res) {
    this.count = res.count;
    this.contribs = res.results;
    this.next = res.next;
    this.prev = res.previous;
    this.emitChange();
  },
  getState: function() {
    return {
      count: this.count,
      contribs: this.contribs,
      next: this.next,
      prev: this.prev
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.count = state.count;
    this.contribs = state.contribs;
    this.next = state.next;
    this.prev = state.prev;
  }
});

