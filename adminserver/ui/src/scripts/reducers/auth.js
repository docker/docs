'use strict';

import assign from 'object-assign';
import consts from 'consts';

const {
    auth,
    loading
} = consts;

const actions = {
  [auth.NOT_SIGNED_IN]: (state) => {
    return assign({}, state, {
      signedIn: false
    });
  },

  [auth.AUTH]: (state, action) => {
    if (!action.ready) {
      return assign({}, state, {
        status: loading.PENDING,
        signedIn: false
      });
    }
    if (action.error) {
      return assign({}, state, {
        status: loading.FAILURE,
        signedIn: false
      });
    }

    // TODO: maybe move this into components
    window.location.reload();
    // Should never reach here
    return assign({}, state, {
      status: loading.SUCCESS,
      signedIn: true
    });
  }
};

export default function(state = {signedIn: true}, data) {
  if (typeof actions[data.type] === 'function') {
    return actions[data.type](state, data);
  }
  return state;
}
