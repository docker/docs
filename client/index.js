import 'babel-polyfill';

import React from 'react';
import { render } from 'react-dom';
import { Provider } from 'react-redux';
import { applyRouterMiddleware, Router, browserHistory } from 'react-router';
import useScroll from 'react-router-scroll';
import configureRoutes from 'routes';
import configureStore from 'store';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import { internalRouterReady } from 'actions/internal';

// Add global css styles
import 'lib/css/global.css';

const initialState = window.__INITIAL_STATE__;
const store = configureStore(initialState);
const routes = configureRoutes(store);


//  Some components use react-tap-event-plugin to listen for touch events.
//  This dependency is temporary and will go away once react v1.0 is released.
//  Until then, be sure to inject this plugin at the start of your app.
//  http://material-ui.com/#/get-started/installation
require('react-tap-event-plugin')();

browserHistory.listen(() => {
  analytics.page();
});

const muiTheme = getMuiTheme();

// useScroll ensures that the page scrolls to the top when you load a new route
// scroll behavior can be extended with a custom shouldUpdateScroll function
// see https://github.com/taion/react-router-scroll for more info

store.dispatch(internalRouterReady());
render((
  <Provider store={store}>
    <MuiThemeProvider muiTheme={muiTheme}>
      <Router
        history={browserHistory}
        routes={routes}
        render={applyRouterMiddleware(useScroll())}
      />
    </MuiThemeProvider>
  </Provider>
), document.getElementById('app'));
