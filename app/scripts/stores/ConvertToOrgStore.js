'use strict';
import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('stores: ConvertToOrgStore');

export default createStore({
  storeName: 'ConvertToOrgStore',
  handlers: {
    CONVERT_TO_ORG_BAD_REQUEST: '_badRequest',
    UPDATE_TO_ORG_OWNER: '_updateOwner'
  },
  initialize: function() {
    this.convertError = false;
    this.error = {};
    this.newOwner = '';
  },
  _badRequest: function(error) {
    this.convertError = true;
    this.error = error;
    this.emitChange();
  },
  _updateOwner: function(payload) {
    this.newOwner = payload.newOwner;
    this.convertError = false;
    this.emitChange();
  },
  getState: function() {
    return {
      convertError: this.convertError,
      error: this.error,
      newOwner: this.newOwner
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.convertError = state.convertError;
    this.error = state.error;
    this.newOwner = state.newOwner;
  }
});
