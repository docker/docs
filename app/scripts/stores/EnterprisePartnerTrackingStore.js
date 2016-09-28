'use strict';

import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('EnterprisePartnerTrackingStore');

export default createStore({
  storeName: 'EnterprisePartnerTrackingStore',
  handlers: {
    ENTERPRISE_PARTNER_RECEIVE_CODE: '_receivePartnerTrackingCode'
  },
  initialize() {
    this.partnervalue = '';
  },
  _receivePartnerTrackingCode({ code }) {
    this.partnervalue = code;
    this.emitChange();
  },
  getState() {
    return {
      partnervalue: this.partnervalue
    };
  },
  dehydrate() {
    return this.getState();
  },
  rehydrate(state) {
    this.partnervalue = state.partnervalue;
  }
});
