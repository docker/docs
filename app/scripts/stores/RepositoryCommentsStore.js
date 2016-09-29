'use strict';
import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('RepositoryCommentsStore');

var RepoCommentsStore = createStore({
  storeName: 'RepositoryCommentsStore',
  handlers: {
    RECEIVE_REPO_COMMENTS: '_receiveRepoComments'
  },
  initialize: function() {
    this.results = [];
    this.prev = null;
    this.next = null;
    this.count = 0;
  },
  _receiveRepoComments: function(res) {
    this.results = res.results;
    this.prev = res.previous;
    this.next = res.next;
    this.count = res.count;
    this.emitChange();
  },
  getState: function() {
    return {
      results: this.results,
      prev: this.prev,
      next: this.next,
      count: this.count
    };
  },
  rehydrate: function(state) {
    this.results = state.results;
    this.prev = state.prev;
    this.next = state.next;
    this.count = state.count;
  },
  dehydrate: function() {
    return this.getState();
  }
});

module.exports = RepoCommentsStore;
