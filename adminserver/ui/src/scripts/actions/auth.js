'use strict';

import { Auth } from 'dtr-js-sdk';
import consts from 'consts';

export function showingAuthForm() {
  return { type: consts.auth.NOT_SIGNED_IN };
}

export function logIn(credentials) {
  return {
    type: consts.auth.AUTH,
    meta: {
      promise: Auth.logIn(credentials)
    }
  };
}

export function logOut() {
  return () => {
    window.location.href = '/logout';
  };
}
