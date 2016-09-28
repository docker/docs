'use strict';

import consts from 'consts';
import assign from 'object-assign';

const {
    logs,
    loading
} = consts;

const actions = {
  [logs.FETCH_LOGS]: (state, action) => {
    if (!action.ready) {
      return assign({}, state, {status: loading.PENDING});
    }

    if (action.error) {
      return assign({}, state, {
        status: loading.FAILURE
      });
    }

    return assign({}, state, {
      status: loading.SUCCESS,
      lines: action.payload.lines
    });
  },

  [logs.OBSERVE_LOG_DATA]: (state, payload) => {
    const lines = payload.data.lines.concat(state.lines);
    return assign({}, state, {
      status: loading.SUCCESS,
      lines
    });
  }
};

export default function logStore(state = {}, payload) {
  if (typeof actions[payload.type] === 'function') {
    return actions[payload.type](state, payload);
  }
  return state;
}
