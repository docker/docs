'use strict';

import { Auth } from 'hub-js-sdk';
import type FluxibleActionContext from '../../../flow-libs/fluxible';
var debug = require('debug')('hub:actions:attemptSignup');

type SignupPayload = {
  username: String;
  password: String;
  email: String;
};

module.exports = function(actionContext: FluxibleActionContext,
                          payload: SignupPayload,
                          done: Function) {
  actionContext.dispatch('SIGNUP_ATTEMPT_START');
  Auth.signup(payload, function(err, res) {
    if (err) {
      debug('error', err);
      if (res.badRequest){
        actionContext.dispatch('SIGNUP_BAD_REQUEST', res.body);
      } else {
        done();
      }
      actionContext.dispatch('SIGNUP_CLEAR_PASSWORD');
    } else {
      actionContext.dispatch('SIGNUP_SUCCESS');
    }
  });
};
