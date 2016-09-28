'use strict';

import { Map, fromJS } from 'immutable';
import consts from 'consts';

const defaultState = new Map({
  // These maps will be raw normalizr results from normalizeSearchResults:
  // {
  //   entities: {
  //     orgs: [...],
  //     users: [...],
  //     repos: [...]
  //   },
  //   result: [...],
  // }
  autocomplete: new Map(),
  search: new Map()
});

const reducers = {
  [consts.SEARCH]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    return state.set('search', fromJS(action.payload));
  }
};

export default function(state = defaultState, action) {
  if (typeof reducers[action.type] === 'function') {
    return reducers[action.type](state, action);
  }
  return state;
}
