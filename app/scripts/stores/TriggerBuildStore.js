'use strict';

const createStore = require('fluxible/addons/createStore');
const debug = require('debug')('TriggerBuildStore');

var TriggerBuildStore = createStore({
  storeName: 'TriggerBuildStore',
  handlers: {
    AB_TRIGGER_SUCCESS: '_abTriggerSuccess',
    AB_TRIGGER_ERROR: '_abTriggerError'
  },
  initialize: function() {
    this.abtrigger = {
      hasError: false,
      success: false
    };
  },
  _abTriggerClear: function() {
    this.abtrigger = {
      hasError: false,
      success: false
    };
    this.emitChange();
  },
  _abTriggerError: function() {
    this.abtrigger.success = false;
    this.abtrigger.hasError = true;
    setTimeout(this._abTriggerClear.bind(this), 3000);
    this.emitChange();
  },
  _abTriggerSuccess: function() {
    this.abtrigger.success = true;
    this.abtrigger.hasError = false;
    setTimeout(this._abTriggerClear.bind(this), 3000);
    this.emitChange();
  },
  getState: function() {
    return {
      abtrigger: this.abtrigger
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.abtrigger = state.abtrigger;
  }
});

module.exports = TriggerBuildStore;
