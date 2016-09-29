'use strict';

import createStore from 'fluxible/addons/createStore';
import _ from 'lodash';
import { STATUS } from './common/Constants';
var debug = require('debug')('EmailNotifStore:');


export default createStore({
  storeName: 'EmailNotifStore',
  handlers: {
    RECEIVE_NOTIFICATIONS: '_receiveNotifications',
    NOTIF_CHECKBOX_CLICK: '_updateNotifications',
    RESET_EMAIL_NOTIFICATIONS_STORE: '_resetBlankSlate',
    SAVE_NOTIFICATIONS_ERROR: '_saveNotifError',
    SAVE_NOTIFICATIONS_SUCCESS: '_saveNotifSuccess'
  },
  initialize: function() {
    // initialize with data from db
    this.starNotification = false;
    this.imgCommentNotification = false;
    this.autoBuildNotification = false;
    this.starNotificationID = -1;
    this.imgCommentNotificationID = -1;
    this.autoBuildNotificationID = -1;
    this.STATUS = STATUS.DEFAULT;
    this.blankNotificationSlate = {};
  },
  _receiveNotifications: function(notifications) {
    for (var i = 0; i < notifications.length; ++i) {
      switch(notifications[i].notification) {
        case 'new_repo_comment':
          this.imgCommentNotification = true;
          this.imgCommentNotificationID = notifications[i].id;
          break;
        case 'new_repo_star':
          this.starNotification = true;
          this.starNotificationID = notifications[i].id;
          break;
        case 'trusted_build_fail':
          this.autoBuildNotification = true;
          this.autoBuildNotificationID = notifications[i].id;
          break;
      }
    }
    this.blankNotificationSlate = this.getState();
    debug(this.blankNotificationSlate);
    this.attempting = false;
    this.emitChange();
  },
  _resetBlankSlate: function() {
    var slate = this.blankNotificationSlate;
    debug('RESET EMAIL NOTIF BLANK SLATE');
    this.starNotification = slate.starNotification;
    this.imgCommentNotification = slate.imgCommentNotification;
    this.autoBuildNotification = slate.autoBuildNotification;
    this.starNotificationID = slate.starNotificationID;
    this.imgCommentNotificationID = slate.imgCommentNotificationID;
    this.autoBuildNotificationID = slate.autoBuildNotificationID;
    this.STATUS = STATUS.DEFAULT;
    this.emitChange();
  },
  _updateNotifications: function(cboxType) {
    this.STATUS = STATUS.DEFAULT;
    switch(cboxType) {
      case 'starNotification':
        this.starNotification = !this.starNotification;
        break;
      case 'imgCommentNotification':
        this.imgCommentNotification = !this.imgCommentNotification;
        break;
      case 'autoBuildNotification':
        this.autoBuildNotification = !this.autoBuildNotification;
        break;
    }
    this.emitChange();
  },
  _saveNotifError: function() {
    this.STATUS = 'ERROR';
    this.emitChange();
  },
  _saveNotifSuccess: function() {
    this.STATUS = STATUS.SUCCESSFUL;
    this.emitChange();
  },
  hasChanged: function(type) {
    switch (type) {
      case 'auto':
        return (this.autoBuildNotification !== this.blankNotificationSlate.autoBuildNotification);
      case 'star':
        return (this.starNotification !== this.blankNotificationSlate.starNotificationID);
      case 'comment':
        return (this.imgCommentNotification !== this.blankNotificationSlate.imgCommentNotificationID);
      default:
        break;
    }
  },
  getAttempt: function() {
    return this.attempting;
  },
  setAttempt: function(flag) {
    this.attempting = flag;
  },
  getState: function() {
    return {
      starNotification: this.starNotification,
      imgCommentNotification: this.imgCommentNotification,
      autoBuildNotification: this.autoBuildNotification,
      starNotificationID: this.starNotificationID,
      imgCommentNotificationID: this.imgCommentNotificationID,
      autoBuildNotificationID: this.autoBuildNotificationID,
      STATUS: this.STATUS
    };
  },
  dehydrate: function() {
    return _.merge({}, this.getState(), {blankNotificationSlate: this.blankNotificationSlate});
  },
  rehydrate: function(state) {
    this.starNotification = state.starNotification;
    this.imgCommentNotification = state.imgCommentNotification;
    this.autoBuildNotification = state.autoBuildNotification;
    this.starNotificationID = state.starNotificationID;
    this.imgCommentNotificationID = state.imgCommentNotificationID;
    this.autoBuildNotificationID = state.autoBuildNotificationID;
    this.STATUS = state.STATUS;
    this.blankNotificationSlate = state.blankNotificationSlate;
  }
});
