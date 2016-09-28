'use strict';

/**
 * displays alert-style dismissable notifications to the user
 */

const createStore = require('fluxible/addons/createStore');
const debug = require('debug')('NotifyStore');

var NotifyStore = createStore({
  storeName: 'NotifyStore',
  handlers: {
    NEW_ALERT: '_newAlert',
    EXPIRE_ALERT: '_expireAlert',
    EXPIRED_SIGNATURE: '_newDetailAlert'
  },
  initialize: function() {
    // alerts have a timestamp-based key
    this.alerts = {};
  },
  _newAlert: function(obj) {
    this.alerts[+new Date()] = obj.msg;
    debug(this.alerts);
    this.emitChange();
  },
  _newDetailAlert: function(msg) {
    debug(msg);
    this._newAlert({
      msg: 'You have been logged out because your token has expired or was invalid'
    });
  },
  getState: function() {
    return this.state;
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    debug(state);
    this.alerts = state.alerts;
  }
});

module.exports = NotifyStore;
