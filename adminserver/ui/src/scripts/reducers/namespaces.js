'use strict';

import { Map } from 'immutable';
import consts from 'consts';

// For now this only allows us to store search results for a namespace.
const defaultState = new Map({
  // These maps will be raw noramlizr results from normalizeSearchResults:
  // {
  //   entities: {
  //     orgs: [...],
  //     users: [...],
  //     repos: [...]
  //   },
  //   result: [...],
  // }
  search: new Map()
});

const reducers = {
  [consts.SEARCH_NAMESPACES]: (state, action) => {
    if (!action.ready || action.error) {
      return state;
    }
    return state.set('search', action.payload);
  }
};

export default function(state = defaultState, action) {
  if (typeof reducers[action.type] === 'function') {
    return reducers[action.type](state, action);
  }
  return state;
}
