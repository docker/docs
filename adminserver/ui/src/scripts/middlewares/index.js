'use strict';

import thunk from 'redux-thunk';

import logger from './logger.js';
import notificationMiddleware from './notification.js';
import promiseMiddleware from './promise.js';
import { routerMiddleware } from 'react-router-redux';
import { browserHistory } from 'react-router';

// Note: Order matters so export as an array
export const middlewares = [
  thunk,
  notificationMiddleware,
  promiseMiddleware,
  routerMiddleware(browserHistory),
  logger
];
