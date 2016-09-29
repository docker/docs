'use strict';

import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('AccountSettingsLicenseStore');
import _ from 'lodash';

export default createStore({
  storeName: 'AccountSettingsLicensesStore',
  handlers: {
    RECEIVE_LICENSES: '_receiveLicenses'
  },
  initialize: function() {
    this.licenses = [];
    this.attempting = true;
  },
  _receiveLicenses: function(licenses) {
    this.licenses = _.flatten(licenses);
    this.attempting = false;
    this.emitChange();
  },
  getAttempt: function() {
    return this.attempting;
  },
  setAttempt: function(flag) {
    this.attempting = flag;
  },
  getState: function() {
    return {
      licenses: this.licenses,
      attempting: this.attempting
    };
  },
  rehydrate: function(state) {
    this.licenses = state.licenses;
    this.attempting = state.attempting;
  },
  dehydrate: function() {
    return this.getState();
  }
});

