'use strict';
import createStore from 'fluxible/addons/createStore';
const debug = require('debug')('stores: ChangePasswordStore');

var ChangePasswordStore = createStore({
  storeName: 'ChangePasswordStore',
  handlers: {
    CHANGE_PASS_UPDATE: '_updateStore',
    CHANGE_PASS_SUCCESS: '_changePassSuccess',
    CHANGE_PASS_CLEAR: '_clearStore',
    RESET_PASSWORD_SUCCESSFUL: '_changePassSuccess',
    RESET_PASSWORD_ERROR: '_changePassError'
  },
  initialize: function() {
    this.oldpass = '';
    this.newpass = '';
    this.confpass = '';
    this.reset = false;
    this.hasErr = false;
    this.err = '';
  },
  _updateStore: function(payload) {
    this.reset = false;
    this.hasErr = false;
    this.err = '';
    this.oldpass = payload.oldpass;
    this.newpass = payload.newpass;
    this.confpass = payload.confpass;
    this.emitChange();
  },
  _changePassSuccess: function() {
    this._clearStore();
    this.reset = true;
    this.emitChange();
  },
  _changePassError: function(error) {
    this._clearStore();
    this.hasErr = true;
    this.err = error;
    this.emitChange();
  },
  _clearStore: function() {
    this.oldpass = '';
    this.newpass = '';
    this.confpass = '';
    this.reset = false;
    this.hasErr = false;
    this.err = '';
    this.emitChange();
  },
  getState: function() {
    return {
      oldpass: this.oldpass,
      newpass: this.newpass,
      confpass: this.confpass,
      reset: this.reset,
      hasErr: this.hasErr,
      err: this.err
    };
  }

});

module.exports = ChangePasswordStore;
