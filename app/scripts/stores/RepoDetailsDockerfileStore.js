'use strict';

import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('RepoDetailsDockerfileStore');

var RepoDetailsDockerfileStore = createStore({
  storeName: 'RepoDetailsDockerfileStore',
  handlers: {
    RECEIVE_DOCKERFILE_FOR_REPOSITORY: '_receiveDockerfile'
  },
  initialize: function() {
    this.dockerfile = '';
  },
  _receiveDockerfile: function(res) {
      debug('dockerfile', res);
    this.dockerfile = res.contents;
    this.emitChange();
  },
  getState: function() {
      return {
          dockerfile: this.dockerfile
      };
  },
  rehydrate: function(state) {
    this.dockerfile = state.dockerfile;
  },
  dehydrate: function() {
    return this.getState();
  }
});

module.exports = RepoDetailsDockerfileStore;
