import { createStore, applyMiddleware, compose } from 'redux';
import thunkMiddleware from 'redux-thunk';
import rootReducer, { DEFAULT_STATE } from 'reducers';
import promiseMiddleware from 'redux-promise-middleware';
import { errorMiddleware, redirectMiddleware } from './middleware';

const isProduction = process.env.NODE_ENV === 'production';

const middleware = [
  thunkMiddleware,
  promiseMiddleware({
    promiseTypeSuffixes: ['REQ', 'ACK', 'ERR'],
  }),
  errorMiddleware,
  redirectMiddleware,
];

export default function (initialState = DEFAULT_STATE) {
  if (!isProduction) {
    const createLogger = require('redux-logger');
    middleware.push(createLogger({
      level: 'info',
      collapsed: true,
    }));
  }

  const withDevTools = !isProduction && typeof window === 'object' &&
    typeof window.devToolsExtension !== 'undefined';

  const store = createStore(
    rootReducer,
    initialState,
    compose(
      applyMiddleware(...middleware),
      withDevTools ? window.devToolsExtension() : f => f,
    ),
  );

  if (module.hot) {
    // Enable Webpack hot module replacement for reducers
    module.hot.accept('reducers', () => {
      // TODO Babel 6
      // const nextRootReducer = require('../reducers').default;
      const nextRootReducer = require('reducers');
      store.replaceReducer(nextRootReducer);
    });
  }

  return store;
}
