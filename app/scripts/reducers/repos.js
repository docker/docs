'use strict';

import immutable from 'immutable';
import {
  RECEIVE_REPO
} from 'reduxConsts.js';

const defaultState = immutable.fromJS(
  (typeof window !== 'undefined' && window.ReduxApp.repos) || {}
);

const reducers = {
  [RECEIVE_REPO]: (state, action) => {
    return state.clear().merge(action.payload);
  }
};

export default function(state = defaultState, action) {
  const { type } = action;
  if (typeof reducers[type] === 'function') {
    return reducers[type](state, action);
  }
  return state;
}
