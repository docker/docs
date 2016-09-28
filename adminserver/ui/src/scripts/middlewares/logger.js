'use strict';

import createLogger from 'redux-logger';
import Immutable from 'immutable';

import consts from 'consts';

const loggerMiddleware = createLogger({
  predicate: () => {
    // Don't log in production
    return process.env.NODE_ENV !== 'production';
  },
  collapsed: (getState, action) => {
    if (action.ready === false) {
      // For promise actions, the pending action isn't as useful
      return true;
    }

    const reduxRouterPrefix = '@@reduxReactRouter/';
    if (action.type.substr(0, reduxRouterPrefix.length) === reduxRouterPrefix) {
      // Router changes are not super important to look at
      return true;
    }

    const notificationConstants = [
        consts.ADD_GROWL_NOTIFICATION,
        consts.REMOVE_GROWL_NOTIFICATION,
        consts.UPDATE_GROWL_NOTIFICATION,
        consts.ADD_BANNER_NOTIFICATION,
        consts.REMOVE_BANNER_NOTIFICATION,
        consts.UPDATE_BANNER_NOTIFICATION
    ];

    if (notificationConstants.includes(action.type)) {
      // Notifications actions are also rarely useful
      return true;
    }
    return false;
  },
  stateTransformer: (state) => {
    let newState = {};
    for (let key of Object.keys(state)) {
      if (Immutable.Iterable.isIterable(state[key])) {
        newState[key] = state[key].toJS();
      } else {
        newState[key] = state[key];
      }
    }
    return newState;
  }
});

export default loggerMiddleware;
