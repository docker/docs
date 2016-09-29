'use strict';
import createStore from 'fluxible/addons/createStore';
import _ from 'lodash';

var AutobuildStore = createStore({
  storeName: 'AutobuildStore',
  handlers: {
    RECEIVE_SOURCE_REPOS: '_receiveSourceRepos',
    RECEIVE_SOURCE_ACCOUNTS: '_receiveSourceAccount'
  },
  initialize: function() {
    this.githubAccount = null;
    this.githubRepos = [];
    this.bitbucketAccount = null;
    this.bitbucketRepos = [];
    this.gitlabAccount = null;
    this.gitlabRepos = [];
  },
  _receiveSourceRepos: function(res) {
    this.githubRepos = res.github.detail ? [] : res.github;
    this.bitbucketRepos = res.bitbucket.detail ? [] : res.bitbucket;
    this.gitlabRepos = res.gitlab.detail ? [] : res.gitlab;
    this.emitChange();
  },
  _receiveSourceAccount: function(res) {
    this.githubAccount = res.github.detail ? null : res.github;
    this.bitbucketAccount = res.bitbucket.detail ? null : res.bitbucket;
    this.gitlabAccount = res.gitlab.detail ? null : res.gitlab;
    this.emitChange();
  },
  getState: function() {
    return {
      githubAccount: this.githubAccount,
      githubRepos: this.githubRepos,
      bitbucketAccount: this.bitbucketAccount,
      bitbucketRepos: this.bitbucketRepos,
      gitlabAccount: this.gitlabAccount,
      gitlabRepos: this.gitlabRepos
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.githubAccount = state.githubAccount;
    this.githubRepos = state.githubRepos;
    this.bitbucketAccount = state.bitbucketAccount;
    this.bitbucketRepos = state.bitbucketRepos;
    this.gitlabAccount = state.gitlabAccount;
    this.gitlabRepos = state.gitlabRepos;
  }
});

module.exports = AutobuildStore;
