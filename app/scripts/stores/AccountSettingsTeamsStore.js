'use strict';

import createStore from 'fluxible/addons/createStore';

export default createStore({
  storeName: 'AccountSettingsTeamsStore',
  handlers: {
    RECEIVE_LICENSES: '_receiveLicenses'
  },
  initialize() {
    this.licenses = [];
  },
  _receiveLicenses(res) {
    this.licenses = res.results;
    this.emitChange();
  },
  getState() {
    return {
      licenses: this.licenses
    };
  },
  rehydrate(state) {
    this.licenses = state.licenses;
  },
  dehydrate() {
    return this.getState();
  }
});

