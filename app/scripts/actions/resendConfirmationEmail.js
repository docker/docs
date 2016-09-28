'use strict';
import {
  Emails
  } from 'hub-js-sdk';
const debug = require('debug')('hub:actions:resendConfirmationEmail');

var resendConfirmationEmail = function(actionContext, {JWT, emailID}) {
  actionContext.dispatch('RESEND_EMAIL_CONFIRMATION_ATTEMPT_START', emailID);
  Emails.resendConfirmationEmail(JWT, emailID, function(err, res) {
    if (err) {
      debug('error', err);
      actionContext.dispatch('RESEND_EMAIL_CONFIRMATION_FAILED', emailID);
    } else {
      actionContext.dispatch('RESEND_EMAIL_CONFIRMATION_SENT', emailID);
    }
    setTimeout(() => {actionContext.dispatch('RESEND_EMAIL_CONFIRMATION_CLEAR', emailID);}, 4000);
  });
};

module.exports = resendConfirmationEmail;
