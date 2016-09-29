/*global window, process, require */
'use strict';

// Needed to shim Object.assign (used by 3rd party libs)
import 'babel-core/polyfill';

import React from 'react';
import { render } from 'react-dom';
import { Router } from 'react-router';
import FluxRouter from './fluxibleRouter';
import createHistory from 'history/lib/createBrowserHistory';
import app from './app';
import navigateAction from './actions/navigate';
import JWTStore from './stores/JWTStore';
import bootstrapCreateElement from './bootstrapCreateElement';

// This renders a Provider for basic flux integration
import { Provider } from 'react-redux';
import reducers from './reducers';
import enhancedCreateStore from './reduxStore';

require('velocity-animate');
require('velocity-animate/velocity.ui');

let oDebug = require('debug');
if (process.env.ENV !== 'production') {
  oDebug.enable('hub:*');
}
const debug = oDebug('hub:client');
const dehydratedState = window.App; // Sent from the server

if (process.env.ENV !== 'production') {
  window.React = React; // For chrome dev tool support
}

let history = createHistory();
let unlisten = history.listen(function (location) {
  debug('unlisten', location.pathname);
});

// Plug a History into the actionContext
let pluginRouter = Router;
app.plug({
  name: 'HistoryPlugin',
  plugContext() {
    return {
      plugActionContext(actionContext) {
        actionContext.history = history;
      }
    };
  }
});

function onUpdate(context) {
  return (state, callback) => {
    debug('at least the second render');
    if (state) {
      if (!state.jwt) {
        let jwtStore = context.getComponentContext().getStore(JWTStore);
        state.jwt = jwtStore.getJWT();
      }
      context.executeAction(navigateAction, state, callback);
    }
  };
}

function renderApp(context) {
  const mountNode = document.getElementById('app');
  const Routes = app.getComponent();

  // context.reduxStore works client-side as we're setting this directly
  // below in rehydrate.
  const jsx = (
    <Provider store={ context.reduxStore }>
      <FluxRouter history={history}
        routes={Routes}
        onUpdate={onUpdate(context)}
        createElement={bootstrapCreateElement(context)}/>
    </Provider>
  );

  return render(jsx, mountNode, () => { debug('React Rendered'); });
}

// The callback is called after the app has rehydrated any plugins;
// our redux plugin neets to create the store itself.
app.rehydrate(dehydratedState, function(err, context) {
  debug('rehydrating app');
  if (err) {
    throw err;
  }

  // Create a new store and save it to Fluxible's app context.
  context.reduxStore = enhancedCreateStore(reducers);

  window.context = context;
  debug('supposedly the first render');
  // Don't call the action on the first render on top of the server rehydration
  // Otherwise there is a race condition where the action gets executed before
  // render has been called, which can cause the checksum to fail.
  renderApp(context);
});
