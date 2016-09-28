'use strict';

// Note: This file should contain all our subtree reducers.
// ./global.js should not be included.

import { combineReducers } from 'redux';

import { reducer as form } from 'redux-form';
import { routerReducer as router } from 'react-router-redux';

import { reducer as ui } from 'redux-ui';
import status from './status';
import logs from './logs';
import settings from './settings';
import metrics from './metrics';
import auth from './auth';
import repositories from './repositories';
import organizations from './organizations';
import notifications from './notifications';
import updates from './updates';
import users from './users';
import teams from './teams';
import search from './search';
import namespaces from './namespaces.js';

export default combineReducers({
  // External node_modules reducers
  form,
  router,
  ui,

  // Custom app reducers
  status,
  logs,
  settings,
  metrics,
  auth,
  repositories,
  organizations,
  notifications,
  updates,
  users,
  teams,
  search,
  namespaces
});
