'use strict';

import createStore from 'fluxible/addons/createStore';
import _ from 'lodash';
import { STATUS } from './common/Constants';
var debug = require('debug')('AccountNotifStore:');

export default createStore({
  storeName: 'OutboundCommunicationStore',
  handlers: {
    RECEIVE_EMAIL_SUBSCRIPTIONS: '_receiveEmailSubscriptions',
    RESET_OUTBOUND_EMAILS_STORE: '_resetBlankSlate',
    SAVE_OUTBOUND_ERROR: '_saveOutboundError',
    SAVE_OUTBOUND_SUCCESS: '_saveOutboundSuccess',
    UPDATE_OUTBOUND: '_updateOutbound',
    UPDATE_BETA_GROUP: '_updateBetaGroup'
  },
  initialize: function() {
    // initialize with data from db
    /*eslint-disable camelcase */
    this.weeklyDigest = {
      subscribed_emails: [],
      unsubscribed_emails: []
    };
    this.digestEmails = [];
    this.betaGroup = {
      subscribed_emails: [],
      unsubscribed_emails: []
    };
    /*eslint-enable camelcase */
    this.betaEmails = [];
    this.STATUS = STATUS.DEFAULT;
    this.blankOutboundSlate = {};
  },
  _receiveEmailSubscriptions: function(payload) {
    this.weeklyDigest = payload.weeklyDigest;
    this.digestEmails = payload.weeklyDigest.subscribed_emails.concat(payload.weeklyDigest.unsubscribed_emails);
    this.betaGroup = payload.betaGroup;
    this.betaEmails = payload.betaGroup.subscribed_emails.concat(payload.betaGroup.unsubscribed_emails);
    this.blankOutboundSlate = this.getState();
    this.emitChange();
  },
  _resetBlankSlate: function() {
    var slate = this.blankOutboundSlate;
    debug('RESET OUTBOUND BLANK SLATE');
    this.weeklyDigest = slate.weeklyDigest;
    this.digestEmails = slate.digestEmails;
    this.betaGroup = slate.betaGroup;
    this.betaEmails = slate.betaEmails;
    this.STATUS = STATUS.DEFAULT;
    this.emitChange();
  },
  _updateOutbound: function(newList) {
    this.STATUS = STATUS.DEFAULT;
    if (newList.list === 'weekly') {
      this.weeklyDigest = newList.data;
    } else if (newList.list === 'beta') {
      this.betaGroup = newList.data;
    }
    this.emitChange();
  },
  _saveOutboundError: function() {
    this.STATUS = 'ERROR';
    this.emitChange();
  },
  _saveOutboundSuccess: function() {
    this.STATUS = STATUS.SUCCESSFUL;
    this.emitChange();
  },
  getState: function() {
    return {
      weeklyDigest: this.weeklyDigest,
      digestEmails: this.digestEmails,
      betaGroup: this.betaGroup,
      betaEmails: this.betaEmails,
      STATUS: this.STATUS
    };
  },
  dehydrate: function() {
    return _.merge({}, this.getState(), {blankOutboundSlate: this.blankOutboundSlate});
  },
  rehydrate: function(state) {
    this.weeklyDigest = state.weeklyDigest;
    this.digestEmails = state.digestEmails;
    this.betaGroup = state.betaGroup;
    this.betaEmails = state.betaEmails;
    this.STATUS = state.STATUS;
    this.blankOutboundSlate = state.blankOutboundSlate;
  }
});
