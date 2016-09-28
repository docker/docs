'use strict';
var debug = require('debug')('hub:server');
var path = require('path');
if(process.env.NEW_RELIC_LICENSE_KEY && process.env.NEW_RELIC_APP_NAME) {
  process.env.NEW_RELIC_NO_CONFIG_FILE = true;
  require('newrelic');
}
var request = require('superagent');
var bugsnag = require('bugsnag');

import pick from 'lodash/object/pick';
import merge from 'lodash/object/merge';
import React from 'react';
import { match } from 'react-router';
import RoutingContext from 'react-router/lib/RoutingContext';
import bootstrapCreateElement from './bootstrapCreateElement';
import app from './app';
import cookieParser from 'cookie-parser';
import bodyParser from 'body-parser';
import express from 'express';
import favicon from 'serve-favicon';
import navigateAction from './actions/navigate';
import serialize from 'serialize-javascript';
import HtmlComponent from './components/Html';

import { Provider } from 'react-redux';
import reducers from './reducers';
import enhancedCreateStore from './reduxStore';

const server = express();

// don't broadcast we are using express
server.disable('x-powered-by');

// add bugsnag for asynch errors
if (process.env.BUGSNAG_API_KEY) {
  bugsnag.register(process.env.BUGSNAG_API_KEY);
  server.use(bugsnag.requestHandler);
}

server.use(favicon('./favicon.ico'));
server.use('/public', express.static('./public'));

// Add a trailing '/' to the path if there is none
server.use(function(req, res, next) {
  if (req.path.substr(-1) !== '/' && req.path.length > 1) {
    const query = req.url.slice(req.path.length);
    res.redirect(301, req.path + '/' + query);
  } else {
    next();
  }
});

// standard health check endpoint
server.get('/_health', function (req, res) {
  res.send('OK');
});

(function(){
  const redirectToDockerPricing = function(req, res) {
    res.redirect(301, 'https://www.docker.com/pricing');
  };

  const redirectTrialToDockerStore = function(req, res) {
    res.redirect(301, 'https://store.docker.com/bundles/docker-datacenter/purchase?plan=free-trial');
  };

  server.get('/enterprise/trial/', redirectTrialToDockerStore); // redirect DDC Trial page to new Store page
  server.get('/enterprise/', redirectToDockerPricing);
  server.get('/subscriptions/', redirectToDockerPricing);
})();

server.get('/account/signup/', function(req, res, next) {
  res.redirect('/');
});

server.get('/account/forgot-password/', function(req, res, next) {
  res.redirect('/reset-password/');
});

server.get('/account/login/', function(req, res, next) {
  res.redirect('/login/');
});

server.get('/_/', function(req, res, next) {
  res.redirect('/explore/');
});

server.get('/official/', function(req, res, next) {
  res.redirect('/explore/');
});

server.get('/account/accounts/', function(req, res, next) {
  res.redirect('/account/authorized-services/');
});

server.get('/plans/', function(req, res, next) {
  res.redirect('https://www.docker.com/pricing');
});

server.get('/resend-email-confirmation/', function(req, res, next) {
  res.redirect('/reset-password/');
});

//There are two cases now:
//Case 1: No query parameter, just the token | Default activation, with confirmation_key sent to the `activate` endpoint
//Case 2: A `ref` query parameter is sent in the email validation URL | For partners, we send both key and ref to the activation endpoint
server.get('/account/confirm-email/:token', function(req, res, next) {
  if(req.params.token) {
    //If on activate, we get a query parameter called `ref` back from the email link, we store it and send it with the POST
    const { ref } = req.query;
    const { token } = req.params;
    var activateRequestBody = { confirmation_key: token };
    if (ref) {
      activateRequestBody.ref = ref;
    }
    request.post(`${process.env.REGISTRY_API_BASE_URL}/v2/users/activate/`)
           .accept('application/json')
           .send(activateRequestBody)
           .end((err, apiRes) => {
             if (err) {
               debug('sign up error', err);
               //Redirect to Login page for any error.
               //We do not have generic error pages for 400s or 500s.
               res.redirect('/login/');
             } else if (!apiRes || !apiRes.body) {
                debug('api response is empty');
                //Redirect to Login page, when there is no response
                //This is a care case that needs to be handled to make sure
                //that it doesn't crash. See HUB-2094 for further details.
                res.redirect('/login/');
             } else {
               const { redirect_url } = apiRes.body;
               if (redirect_url) {
                 //Redirect to the redirect URL returned by the API
                 res.redirect(redirect_url);
               } else {
                 //Fallback redirect to Login page, when there is no redirect URL
                 res.redirect('/login/');
               }
             }
           });
  } else {
    res.redirect('/login/');
  }
});

server.use(cookieParser());
// 30 days in ms: 2592000000
const expiry = 1000 * 60 * 60 * 24 * 30; // ms * s * m * h * days
const cookieOpts = {
  domain: process.env.COOKIE_DOMAIN,
  httpOnly: true,
  secure: true,
  maxAge: expiry,
  expires: new Date(Date.now() + expiry)
};
server.post('/attempt-login/',
            bodyParser.json(),
            function(req, res, next) {
  res.cookie('token', req.body.jwt, cookieOpts);
  res.end();
});

const isLoggedIn = function(req) {
  return !!(req.cookies && (req.cookies.token || req.cookies.jwt));
};

server.get('/account/billing-plans/', function(req, res, next) {
  if (!isLoggedIn(req)) {
    return res.redirect('/billing-plans/');
  }
  next();
});

server.use(function(req, res, next) {
  if (req.method !== 'GET') {
    return next();
  }
  if (isLoggedIn(req) && (['/login/', '/reset-password/', '/register/'].indexOf(req.path) !== -1 || req.path.indexOf('/account/password-reset-confirm') === 0)) {
    res.redirect('/');
  } else if (!isLoggedIn(req) && req.path.indexOf('password-reset-confirm') === -1 && req.path.indexOf('/account/') === 0) {
    res.redirect('/login/');
  } else {
    next();
  }
});

server.post('/attempt-logout/', function(req, res, next) {
  /**
   * Delete the old cookie when we see it on logout.
   */
  const oldCookieOpts = merge({},
                              cookieOpts,
                              {
                                domain: '.docker.com'
                              });
  res.clearCookie('jwt', oldCookieOpts);
  res.clearCookie('token', cookieOpts);
  res.end();
});

server.post('/oauth/github-attempt/',
            bodyParser.json(),
            function(req, res, next) {
  res.cookie('ghOauthKey', req.body.ghk, cookieOpts);
  res.end();
});

server.post('/oauth/github-done/', function(req, res, next) {
  res.clearCookie('ghOauthKey', cookieOpts);
  res.end();
});

server.use(function(req, res, next) {
  // We may need to whitelist OPTIONS
  if (req.method !== 'GET') {
    res.end('This server does not respond to non-GET requests');
  } else {
    next();
  }
});

server.use(function(req, res, next) {
  // Within each request create a new Redux store from all of our reducers
  // so that state is unique per request.
  const store = enhancedCreateStore(reducers);
  const context = app.createContext({
    reduxStore: store
  });

  debug('context:', context, context.reduxStore);

  //We get the Routes that have been created in the FluxibleComponent
  const routes = app.getComponent();

  const originalURL = req.originalUrl;
  //We use the 'match' API to match the created routes with the current location (req.originalURL)
  debug('matching route', originalURL);
  match({ routes, location: originalURL }, (routerError, redirectLocation, renderProps) => {
    // match uses createRoutes for history
    //TODO: handle redirect, not found and errors
    //TODO: need to handle generic 404s, 500s, 301s
    //if (redirectLocation) {
    //TODO: redirects need to be handled here
    //  res.redirect(301, redirectLocation.pathname + redirectLocation.search);
    //}
    //else if (error) {
    //TODO: Render a nice 500 page with error displayed | HOPE THIS NEVER HAPPENS
    //  res.send(500, error.message);
    //}
    //else if (renderProps == null) {
    //TODO: Probably render the 404 page here
    //  res.send(404, 'Not found');
    //}

    //If router errors out, bail
    if (routerError) {
      debug('Error in the Router', routerError);
      res.end(routerError);
    }
    // whitelist cookies from express into renderProps
    if (req.cookies) {
      renderProps.cookies = pick(req.cookies, ['token', 'ghOauthKey']);
      // For backward compat since we changed the cookie name
      renderProps.cookies.jwt = renderProps.cookies.token;
    }

    // Set the props, so the server knows if the user is logged in
    if (renderProps.cookies.jwt) {
      renderProps.JWT = renderProps.cookies.jwt;
    }

    /**
     * Execute navigate action to load data (we block the render until the data is
     * completely loaded)
     * You can see the actual server side render happens only after the
     * `navigateAction` calls `done()` somewhere
     */
    context.executeAction(navigateAction, renderProps, function() {
      debug('Exposing context state', context);
      debug('EXPOSING RENDER PROPS', renderProps);
      let serializedApp;
      let reduxApp;
      try {
        /*
        NOTE: If we have any html or request responses saved in the store
          - serialize will not be able to parse this and will crash the node server
        */
        serializedApp = serialize(app.dehydrate(context));
        reduxApp = serialize(store.getState());
      } catch (err) {
        debug('SERIALIZATION FAILURE: ', err);
      }
      const exposed = `window.App=${serializedApp}; window.ReduxApp = ${reduxApp};`;
      debug('Rendering Application component into html');

      // This is the Router 1.0.0 recommended way of doing server side rendering
      // Also add a Provider around the routingContext for Redux.
      // NOTE: We're defining our redux store above directly within the app context
      const routingContext = (
        <Provider store={ store }>
          <RoutingContext {...renderProps}
            createElement={bootstrapCreateElement(context)}/>
        </Provider>
      );

      debug('rendering html');
      var html = React.renderToStaticMarkup(
        <HtmlComponent
          state={ exposed }
          markup={ React.renderToString(routingContext) }/>
      );

      res.send(html);
    });
  });
});

// add bugsnag for error handling middleware
if(process.env.BUGSNAG_API_KEY) {
  server.use(bugsnag.errorHandler);
}

// add generic error catching middleware so the server doesn't crash
server.use(function catchError(err, req, res, next) {
  const message = err.stack ? err.stack.replace(/\n/g, '') : '';
  const errorLog = {
    time: (new Date()).toISOString(),
    service: 'hub-web-v2',
    message
  };
  console.error(errorLog); // eslint-disable-line no-console
});

const port = process.env.PORT || 3000;

// Stop the server if the process terminates
const runningServer = server.listen(port, function onListen() {
  process.on('exit', runningServer.close.bind(runningServer));
  debug('Listening on port ' + port);
});
