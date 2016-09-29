'use strict';

import immutable, { Map } from 'immutable';
import endsWith from 'lodash/string/endsWith';
import isArray from 'lodash/lang/isArray';
import { ERROR } from 'reduxConsts.js';
import { StatusRecord } from 'records';

const defaultState = immutable.fromJS(
  (typeof window !== 'undefined' && window.ReduxApp.status) || {}
);

// This reducer listens for status updates from the SDK middleware
// and automatically stores the status within the `statusKey`.
//
// Example: If statusKey = ['deleteRepoTag', 'latest'] and status = 'ATTEMPTING',
// then state.status.deleteRepoTag.latest would be `ATTEMPTING`.
export default function(state = defaultState, action) {
  // The status reducer only acts on actions ending in _STATUS.
  // Ignore anything else and return the default state.
  if (!endsWith(action.type, `_STATUS`)) {
    return state;
  }

  const { statusKey, status, data } = action.payload;
  // We're using setIn, so if the statusKey is a string it needs
  // to be wrapped in an array.
  const sk = isArray(statusKey) ? statusKey : [statusKey];

  if (status === ERROR) {
    // Store the status and error response from the API within
    // the state.
    return state.setIn(sk, new StatusRecord({ status, error: data }));
  }

  // On ATTEMPTING or SUCCESS we only want to store the status;
  // storing action.payload.data would store the entire API response
  // which our other reducers should be handling.
  return state.setIn(sk, new StatusRecord({ status }));
}
