'use strict';

import immutable, { Map } from 'immutable';
import {
  RECEIVE_NAUTILUS_TAGS_FOR_REPOSITORY,
  RECEIVE_TAGS_FOR_REPOSITORY,
  DELETE_REPO_TAG,
  SUCCESS
} from 'reduxConsts.js';

// Use the serialized redux data from the universal app loading if it exists.
//
// Shape of state:
// {
//   'namespace': {
//     'reponame': {
//         tags: {
//           'latest': { ...data },
//           '14.04': { ...data },
//           ...
//         },
//         result: [1, 2, 3, ...] // Array of repo IDs as ordered by hub API
//     },
//     ...
//   }
// }
//
// We use a nested map of namespace:reponame keys to a list of tags to ensure
// that we can merge the nautilus and hub API responses together without
// clearing inbetween.
const defaultState = immutable.fromJS(
  (typeof window !== 'undefined' && window.ReduxApp.tags) || {}
);

// mergeTagsIntoState accepts a namespace, reponame and normalized tag
// information and merges them into the given state.
//
// This is used when tags from the hub and nautilus API are loaded.
const mergeTagsIntoState = (state, action) => {
  const { namespace, reponame, tags } = action.payload;
  const path = [namespace, reponame, 'tags'];
  const { tag } = tags.entities;

  return state.setIn(
    path,
    // Get the existing tags for this namespace/repo and merge the normalized
    // tags recursively.  If the namespace/repo pair doesn't exist this returns
    // a new map.
    //
    // Merge function ensures that in the event of a conflict where the new value
    // is undefined or null, it does not overwrite the existing value
    state.getIn(path, new Map()).mergeDeep(tag)
  );
};

const maybeDeleteTag = (state, action) => {
  if (action.payload.status === SUCCESS) {
    // Remove this tag from our reducer.
    const { namespace, reponame, tagName } = action.payload;
    return state.withMutations(s => {
      let result = s.getIn([namespace, reponame, 'result']);
      result = result.filter(tag => tag !== tagName);
      s.deleteIn([namespace, reponame, 'tags', tagName]);
      s.setIn([namespace, reponame, 'result'], result);
      return s;
    });
  }
  return state;
};

const reducers = {
  [RECEIVE_TAGS_FOR_REPOSITORY]: (state, action) => {
    // Add the result array of ordered tags from the hub API response,
    // then merge tags in
    const { namespace, reponame, tags } = action.payload;
    const path = [namespace, reponame, 'result'];
    const { result } = tags;
    state = state.setIn(path, result);
    return mergeTagsIntoState(state, action);
  },
  [RECEIVE_NAUTILUS_TAGS_FOR_REPOSITORY]: mergeTagsIntoState,
  [`${DELETE_REPO_TAG}_STATUS`]: maybeDeleteTag
};

export default function(state = defaultState, action) {
  const { type } = action;
  if (typeof reducers[type] === 'function') {
    return reducers[type](state, action);
  }
  return state;
}
