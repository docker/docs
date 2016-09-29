'use strict';
import createStore from 'fluxible/addons/createStore';
import _ from 'lodash';

var AutobuildSourceRepositoriesStore = createStore({
  storeName: 'AutobuildSourceRepositoriesStore',
  handlers: {
    RECEIVE_LINKED_REPO_SOURCES: '_receiveLinkedRepos',
    LINKED_REPO_SOURCES_ERROR: '_linkedReposError',
    SET_LINKED_REPO_TYPE: '_setType'
  },
  initialize: function() {
    this.repos = [];
    this.type = '';
    this.error = '';
  },
  _setType: function(type) {
    this.type = type;
    this.emitChange();
  },
  _receiveLinkedRepos: function(linkedRepos) {
    this.repos = linkedRepos;
    this.emitChange();
  },
  _linkedReposError: function(err) {
    this.error = 'Please check if you have any repositories setup on ' + this.type + '.';
    this.emitChange();
  },
  getState: function() {
    return {
      repos: this.repos,
      type: this.type,
      error: this.error
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.repos = state.repos;
    this.type = state.type;
    this.error = state.error;
  }
});

module.exports = AutobuildSourceRepositoriesStore;
