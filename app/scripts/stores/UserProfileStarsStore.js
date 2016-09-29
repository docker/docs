'use strict';

import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('stores: UserProfileStarsStore');

export default createStore({
  storeName: 'UserProfileStarsStore',
  handlers: {
    RECEIVE_PROFILE_STARRED_REPOS: '_receiveStarredRepos'
  },
  initialize() {
    this.repos = [];
    this.next = null;
    this.prev = null;
  },
  _receiveStarredRepos(res) {
    this.starred = res.results;
    this.next = res.next;
    this.prev = res.previous;
    this.emitChange();
  },
  getState() {
    return {
      starred: this.starred,
      next: this.next,
      prev: this.prev
    };
  },
  dehydrate() {
    return this.getState();
  },
  rehydrate(state) {
    this.starred = state.starred;
    this.next = state.next;
    this.prev = state.prev;
  }
});
