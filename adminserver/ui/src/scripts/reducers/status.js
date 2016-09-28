'use strict';

import { Map } from 'immutable';
import consts from 'consts';

const defaultState = new Map();

const actions = {
  [consts.status.RESET_STATUS]: () => {
    return defaultState;
  }
};

/**
 * Use in conjunction with promise middleware.
 */
export default function statusStore(state = new Map(), action) {
  if (typeof actions[action.type] === 'function') {
    return actions[action.type](state, action);
  }

  if (action.meta && action.meta.promiseId) {
    if (!action.type) {
      throw new Error('Undefined action type');
    }

    const keyPath = [action.type, ...action.meta.promiseScope];

    if (!action.ready) {
      return state.setIn(keyPath, new Map({
        promiseId: action.meta.promiseId,
        status: consts.loading.PENDING
      }));
    }

    if (action.error) {
      return state.setIn(keyPath, new Map({
        status: consts.loading.FAILURE,
        error: action.error
      }));
    }

    return state.setIn(keyPath, new Map({ status: consts.loading.SUCCESS }));
  }
  return state;
}
