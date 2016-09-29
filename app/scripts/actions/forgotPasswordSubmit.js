/* @flow */
'use strict';
import {
  Users
  } from 'hub-js-sdk';
const debug = require('debug')('hub:actions:forgotPasswordSubmit');

var forgotPasswordSubmit = function(actionContext:{ dispatch : Function }, { email }) {
  Users.forgotPassword(email, function(err, res) {
    if (err) {
      debug('forgotPassword error', err);
    } else {
      actionContext.dispatch('FORGOT_PASSWORD_SENT', res.body);
    }
  });
};

module.exports = forgotPasswordSubmit;
