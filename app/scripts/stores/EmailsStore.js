'use strict';

import createStore from 'fluxible/addons/createStore';
import merge from 'lodash/object/merge';
import has from 'lodash/object/has';
import find from 'lodash/collection/find';
import cloneDeep from 'lodash/lang/cloneDeep';
import { STATUS, EMAILSTATUS } from './emailsstore/Constants';

var debug = require('debug')('EmailsStore');

var EmailsStore = createStore({
  storeName: 'EmailsStore',
  handlers: {
    RECEIVE_EMAILS: '_receiveEmails',
    CHANGE_ROUTE: '_resetState',
    ADD_EMAIL_INVALID: '_addEmailInvalid',
    ADD_EMAIL_SUCCESS: '_addEmailSuccess',
    START_SAVE_ACTION: '_startSaveAction',
    FINISH_SAVE_ACTION: '_finishSaveAction',
    RESEND_EMAIL_CONFIRMATION_ATTEMPT_START: '_resendEmailConfirmationAttemptStart',
    RESEND_EMAIL_CONFIRMATION_SENT: '_resendEmailConfirmationSent',
    RESEND_EMAIL_CONFIRMATION_FAILED: '_resendEmailConfirmationFail',
    RESEND_EMAIL_CONFIRMATION_CLEAR: '_resendClear',
    UPDATE_ADD_EMAIL: '_updateAddEmail'
  },
  initialize: function() {
    this.STATUS = STATUS.DEFAULT;

    this._cleanSlate = {
      emails: []
    };
    this.emails = [];
    this.emailConfirmations = {};
    this.addEmail = '';
    this.addError = '';
  },
  _startSaveAction() {
    this.STATUS = STATUS.SAVING;
    this.emitChange();
  },
  _finishSaveAction() {
    this.STATUS = STATUS.DEFAULT;
    this.emitChange();
  },
  _resendEmailConfirmationAttemptStart(emailID) {
    debug('RESEND CONFIRMATION START');
    this.emailConfirmations = merge(this.emailConfirmations, {
      [emailID]: EMAILSTATUS.ATTEMPTING
    });
    this.emitChange();
  },
  _resendEmailConfirmationSent(emailID) {
    debug('RESEND CONFIRMATION SENT');
    this.emailConfirmations = merge(this.emailConfirmations, {
      [emailID]: EMAILSTATUS.SUCCESS
    });
    this.emitChange();
  },
  _resendEmailConfirmationFail(emailID) {
    debug('RESEND CONFIRMATION FAIL');
    this.emailConfirmations = merge(this.emailConfirmations, {
      [emailID]: EMAILSTATUS.FAILED
    });
    this.emitChange();
  },
  _resendClear(emailID) {
    debug('RESEND CONFIRMATION CLEAR');
    this.emailConfirmations = merge(this.emailConfirmations, {
      [emailID]: ''
    });
    this.emitChange();
  },
  _receiveEmails: function(payload) {
    debug(payload);
    this.initialize();
    this._cleanSlate = {
      /**
       * We cloneDeep here because otherwise this.emails
       * and this._cleanSlate.emails will refer to the
       * same array causing unintuitive behavior.
       */
      emails: cloneDeep(payload.emails)
    };
    this.emails = payload.emails;
    this.emitChange();
  },
  _addEmailInvalid(error) {
    this.addError = error[0];
    this.emitChange();
  },
  _addEmailSuccess() {
    this.addEmail = '';
    this.emitChange();
  },
  _resetState() {
    var {emails} = this._cleanSlate;
    this.emails = emails.slice();
    this.addError = '';
    this.emitChange();
  },
  _updateAddEmail(email) {
    this.addEmail = email;
    this.emitChange();
  },
  isCleanSlatePrimaryEmail: function(email: string) {
    debug('cleanSlate.emails', this._cleanSlate.emails, email);
    /**
     * A function that answers "is this email address a primary email
     * address?" with respect to the database, not with respect
     * to the state of the client side application
     */
    var primaryEmail = find(this._cleanSlate.emails, function(obj) {
      return obj.email === email && obj.primary === true;
    });

    debug('primaryEmail', !!primaryEmail);
    return !!primaryEmail;
  },
  getCleanSlatePrimaryEmailID() {
    return find(this._cleanSlate.emails, function(obj) {
      return obj.primary === true;
    }).id;
  },
  getEmails: function() {
    return {
      emails: this.emails
    };
  },
  getState: function() {
    return {
      STATUS: this.STATUS,
      emails: this.emails,
      addEmail: this.addEmail,
      addError: this.addError,
      emailConfirmations: this.emailConfirmations
    };
  },
  dehydrate() {
    return merge({},
                 this.getState(),
                 {
                   _cleanSlate: this._cleanSlate
                 });
  },
  rehydrate(state) {
    this._cleanSlate = state._cleanSlate;
    this.addEmail = state.addEmail;
    this.emails = state.emails.slice(0);
    this.emailConfirmations = state.emailConfirmations;
    this.addError = state.addError;
  }
});

module.exports = EmailsStore;
