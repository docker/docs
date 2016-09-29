'use strict';

import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('DashboardReposStore');
import { STATUS } from './common/Constants';
const {
  ATTEMPTING,
  DEFAULT,
  SUCCESSFUL
} = STATUS;

var DashboardReposStore = createStore({
  storeName: 'DashboardReposStore',
  handlers: {
    RECEIVE_REPOS: '_receiveRepos',
    LOGOUT: 'initialize',
    DASHBOARD_REPOS_STORE_ATTEMPTING_GET_REPOS: '_startGetRepos',
    DASHBOARD_REPOS_STORE_ATTEMPTING_GET_ALL_REPOS: '_startGetAllRepos',
    DASHBOARD_REPOS_STORE_RECEIVE_ALL_REPOS_SUCCESS: '_receiveAllRepos'
  },
  initialize() {
    this.repos = [];
    this.count = 0;
    this.next = null;
    this.prev = null;
    this.STATUS = DEFAULT;
  },
  getState() {
    return {
      repos: this.repos,
      count: this.count,
      next: this.next,
      prev: this.prev,
      STATUS: this.STATUS
    };
  },
  _startGetRepos: function() {
    this.STATUS = DEFAULT;
    this.emitChange();
  },
  _startGetAllRepos: function() {
    this.STATUS = ATTEMPTING;
    this.emitChange();
  },
  _receiveRepos(res) {
    debug(res);
    this.repos = res.results;
    this.count = res.count;
    this.next = res.next;
    this.prev = res.previous;
    this.emitChange();
  },
  _receiveAllRepos(res) {
    this.STATUS = SUCCESSFUL;
    this._receiveRepos(res);
  },
  dehydrate() {
    return this.getState();
  },
  rehydrate(state) {
    this.repos = state.repos;
    this.count = state.count;
    this.next = state.next;
    this.prev = state.prev;
    this.STATUS = state.STATUS;
  }
});

module.exports = DashboardReposStore;
