'use strict';

import createStore from 'fluxible/addons/createStore';
import filter from 'lodash/collection/filter';
import { PENDING_DELETE } from 'common/enums/RepoStatus';
var debug = require('debug')('stores: UserProfileStore');

export default createStore({
  storeName: 'UserProfileReposStore',
  handlers: {
    RECEIVE_PROFILE_REPOS: 'receiveRepos'
  },
  initialize() {
    this.repos = [];
    this.next = null;
    this.prev = null;
  },
  removePendingDeleteRepos(repos) {
    //Remove repos that are in pending delete state from user profile repos
    return filter(repos, (repo) => {
      const { status } = repo;
      return status !== PENDING_DELETE;
    });
  },
  receiveRepos(res) {
    this.repos = this.removePendingDeleteRepos(res.results);
    this.next = res.next;
    this.prev = res.previous;
    this.emitChange();
  },
  getState() {
    return {
      repos: this.repos,
      next: this.next,
      prev: this.prev
    };
  },
  dehydrate() {
    return this.getState();
  },
  rehydrate(state) {
    this.repos = state.repos;
    this.next = state.next;
    this.prev = state.prev;
  }
});
