'use strict';

import createStore from'fluxible/addons/createStore';

var ReposStore = createStore({
  storeName: 'RepositoriesPageStore',
  handlers: {
    RECEIVE_REPOS: '_receiveRepos'
  },
  initialize: function() {
    this.repos = [];
    this.previous = null;
    this.next = null;
    this.count = null;
  },
  _receiveRepos: function(res) {
    this.repos = res.results;
    this.previous = res.previous;
    this.next = res.next;
    this.count = res.count;
    this.emitChange();
  },
  getState: function() {
    return {
      repos: this.repos,
      previous: this.previous,
      next: this.next,
      count: this.count
    };
  },
  rehydrate: function(state) {
    this.repos = state.repos;
    this.previous = state.previous;
    this.next = state.next;
    this.count = state.count;
  },
  dehydrate: function() {
    return this.getState();
  }
});

module.exports = ReposStore;
