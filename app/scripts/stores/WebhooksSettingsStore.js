'use strict';
import createStore from 'fluxible/addons/createStore';
const debug = require('debug')('WebhooksSettingsStore');

var WebhooksSettingsStore = createStore({
  storeName: 'WebhooksSettingsStore',
  handlers: {
    RECEIVE_WEBHOOKS: '_receiveWebhooks'
  },
  initialize() {
    this.pipelines = [];
  },
  _receiveWebhooks(payload) {
    debug(payload);
    this.pipelines = payload.results;
    this.emitChange();
  },
  getState() {
    return {
      pipelines: this.pipelines
    };
  },
  dehydrate() {
    return this.getState();
  },
  rehydrate(state) {
    this.pipelines = state.pipelines;
  }
});

module.exports = WebhooksSettingsStore;
