'use strict';
import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('stores: JWTStore');
import cookie from 'cookie';

export default createStore({
  storeName: 'JWTStore',
  handlers: {
    RECEIVE_JWT: '_receiveJWT',
    LOGOUT: '_logout',
    LOGOUT_ERROR: '_logoutError',
    EXPIRED_SIGNATURE: '_setExpiredSignature'
  },
  _receiveJWT(jwt) {
    this.jwt = jwt;
    this.signatureIsExpired = false;
    this.emitChange();
  },
  _logoutError(err) {
    debug(err + ' Logout did not complete cleanly on the server');
    this._logout();  //we logout on the client side anyway
  },
  _logout() {
      this.jwt = null;
      this.emitChange();
  },
  _logoutWithNotification(){
    debug('Logging out due to invalid Signature');
    this._logout();
  },
  _setExpiredSignature(){
    this.signatureIsExpired = true;
    this.jwt = null;
    this.emitChange();
  },
  getJWT() {
    return this.jwt;
  },
  getState() {
    return {
      jwt: this.jwt,
      signatureIsExpired: this.signatureIsExpired
    };
  },
  isLoggedIn() {
    //Return true if user is logged in
    return !!this.jwt;
  },
  dehydrate() {
    if(this.signatureIsExpired) {
      return {
        jwt: null,
        signatureIsExpired: true
      };
    } else {
      return {
        jwt: this.jwt,
        signatureIsExpired: false
      };
    }
  },
  rehydrate(state) {
    debug('rehydrate', state);
    if(state.signatureIsExpired) {
      debug('signatureIsExpired');
      this._logoutWithNotification();
    } else {
      debug('signatureIsValid');
      this.signatureIsExpired = state.signatureIsExpired;
      this._receiveJWT(state.jwt);
    }
  }
});
