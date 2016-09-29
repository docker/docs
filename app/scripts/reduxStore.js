'use strict';

// This creates a new redux store by importing our reducers and middleware
// and combining the two.
import { applyMiddleware, compose, createStore } from 'redux';
import { Iterable } from 'immutable';
import reducers from './reducers';
import sdkMiddleware from './middlewares/sdk.js';
import createLogger from 'redux-logger';
const debug = require('debug')('hub:redux:logger');


// Logger must always be the last middleware in applyMiddleware
const logger = createLogger({
  predicate: () => process.env.ENV === `development`,
  // Use the debug import as our logger
  logger: {log: debug},
  // Transform any immutableJS maps and iterables into their standard JS
  // counterparts. This means you can inspect state within the console.
  stateTransformer: (state) => {
    let newState = {};
    Object.keys(state).forEach(key => {
      newState[key] = state[key];
      if (Iterable.isIterable(state[key])) {
        newState[key] = state[key].toJS();
      }
    });
    return newState;
  }
});

// Compose creates a new function by taking store enhancers (such as middleware
// and any external enhancers) which modifies the createStore function.
const enhancedCreateStore = compose(
  applyMiddleware(sdkMiddleware, logger)
)(createStore);

export default enhancedCreateStore;
