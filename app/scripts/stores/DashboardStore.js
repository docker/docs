'use strict';

import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('DashboardStore');

var DashboardStore = createStore({
  storeName: 'DashboardStore',
  handlers: {
    RECEIVE_STARRED: '_receiveStarredRepos',
    RECEIVE_CONTRIB: '_receiveContribRepos',
    RECEIVE_ACTIVITY_FEED: '_receiveActivityFeed',
    LOGOUT: 'initialize'
  },
  initialize: function() {
    this.user = {};
    this.org = '';
    this.starred = [];
    this.contribs = [];
    this.feed = [];
  },
  getInitState: function() {
    return {
      starred: [],
      contribs: [],
      org: '',
      feed: [],
      user: {}
    };
  },
  _receiveStarredRepos: function(repos) {
    this.starred = repos.results;
    this.emitChange();
  },
  _receiveContribRepos: function(repos) {
    this.contribs = repos.results;
    this.emitChange();
  },
  _receiveActivityFeed: function(feed) {
    this.feed = feed;
    this.emitChange();
  },
  getState: function() {
    return {
      starred: this.starred,
      contribs: this.contribs,
      org: this.org,
      feed: this.feed,
      user: this.user
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.starred = state.starred;
    this.contribs = state.contribs;
    this.feed = state.feed;
    this.org = state.org;
    this.user = state.user;
  }
});

module.exports = DashboardStore;
