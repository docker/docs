'use strict';

// This combines all reducers from internal and external packages to create
// a redux store.
import { combineReducers } from 'redux';
import repos from './repos';
import scans from './scans';
import status from './status';
import tags from './tags';
import { reducer as ui } from 'redux-ui';

export default combineReducers({
  // external reducers
  ui,
  // middleware reducers
  // app-specific reducers
  repos,
  scans,
  status,
  tags
});
