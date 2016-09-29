'use strict';

import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('DashboardStarsStore');

export default createStore({
  storeName: 'DashboardStarsStore',
  handlers: {
    RECEIVE_STARRED: '_receiveStarredRepos',
    LOGOUT: 'initialize'
  },
  initialize: function() {
    this.count = 0;
    this.starred = [];
    this.next = null;
    this.prev = null;
  },
  _receiveStarredRepos: function(res) {
    this.count = res.count;
    this.starred = res.results;
    this.next = res.next;
    this.prev = res.previous;
    this.emitChange();
  },
  getState: function() {
    return {
      count: this.count,
      starred: this.starred,
      next: this.next,
      prev: this.prev
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.count = state.count;
    this.starred = state.starred;
    this.next = state.next;
    this.prev = state.prev;
  }
});

