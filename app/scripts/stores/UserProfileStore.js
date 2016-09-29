'use strict';

import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('UserProfileStore');

export default createStore({
  storeName: 'UserProfileStore',
  handlers: {
    RECEIVE_PROFILE_USER: '_receiveUser',
    USER_PROFILE_404: '_fourOHfour'
  },
  initialize() {
    this.STATUS = 'DEFAULT';
    this.user = {};
  },
  _fourOHfour() {
    this.STATUS = '404';
    this.emitChange();
  },
  _receiveUser(user) {
    this.STATUS = 'DEFAULT';
    this.user = user;
    this.emitChange();
  },
  getState() {
    return {
      user: this.user,
      STATUS: this.STATUS
    };
  },
  dehydrate() {
    return this.getState();
  },
  rehydrate(state) {
    this.user = state.user;
    this.STATUS = state.STATUS;
  }
});
