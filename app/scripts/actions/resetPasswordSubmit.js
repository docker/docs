/* @flow */
'use strict';
import {
  Users
  } from 'hub-js-sdk';
const debug = require('debug')('hub:actions:resetPasswordSubmit');

var resetPasswordSubmit = function(actionContext:{ dispatch : Function }, { uidb64, reset_token, password_1, password_2 }) {
  Users.resetPassword(uidb64, password_1, password_2, reset_token, function(err, res) {
    if (err) {
      debug('error', err);
      actionContext.dispatch('RESET_PASSWORD_ERROR', res.body);
    } else if (res.ok) {
      actionContext.dispatch('RESET_PASSWORD_SUCCESSFUL');
      actionContext.history.push('/account/password-reset-confirm/success/');
      actionContext.dispatch('CHANGE_PASS_CLEAR', {});
    }
  });
};

module.exports = resetPasswordSubmit;
