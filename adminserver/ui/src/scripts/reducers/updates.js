'use strict';

import assign from 'object-assign';

import consts from 'consts';

const actions = {
  [consts.settings.UPDATES]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    return assign({}, state, action.payload);
  },
  [consts.updates.UPDATE_AVAILABLE]: (state, data) => {
    return assign({}, state, data.payload);
  }
};

export default function(state = {}, data) {
  if (typeof actions[data.type] === 'function') {
    return actions[data.type](state, data);
  }
  return state;
}
