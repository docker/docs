'use strict';

//modified from nautilus-ui/src/scripts/reducers/scans.js
import immutable from 'immutable';
import { RECEIVE_SCANNED_TAG_DATA } from 'reduxConsts.js';

// Map of entities within each scan
const defaultState = immutable.fromJS(
  (typeof window !== 'undefined' && window.ReduxApp.scans) || {}
);

const reducers = {
  [ RECEIVE_SCANNED_TAG_DATA ]: (state, action) => {
    // Here we only ever save this current scan from the repoDetailsScannedTag
    // action.  This means that our scans reducer only ever has one scan - for
    // the current page.
    return state.clear().merge(action.payload.entities);
  }
};

export default function(state = defaultState, action) {
  if (typeof reducers[action.type] === 'function') {
    return reducers[action.type](state, action);
  }
  return state;
}
