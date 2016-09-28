'use strict';

import React from 'react';
import createHashHistory from 'history/lib/createHashHistory';
import { createRoutes, useRoutes, RoutingContext } from 'react-router';
import { routes } from 'react-router/lib/PropTypes';
var debug = require('debug')('FluxibleRouter');

const { func, object } = React.PropTypes;

/**
 * A <FluxibleRouter> is a high-level API for automatically setting up
 * a router that renders a <RoutingContext> with all the props
 * it needs each time the URL changes.
 */
const FluxibleRouter = React.createClass({

  propTypes: {
    history: object,
    children: routes,
    routes, // alias for children
    createElement: func,
    onError: func,
    onUpdate: func,
    parseQueryString: func,
    stringifyQuery: func
  },

  getInitialState() {
    return {
      firstRender: true,
      location: null,
      routes: null,
      params: null,
      components: null
    };
  },

  handleError(error) {
    if (this.props.onError) {
      this.props.onError.call(this, error);
    } else {
      // Throw errors by default so we don't silently swallow them!
      throw error; // This error probably occurred in getChildRoutes or getComponents.
    }
  },

  //============================================================================
  //ComponentWillMount is the only addition to this Router
  //Needed in order to provide location/pathname to onUpdate from client.js
  /* eslint-disable */
  componentWillMount() {
    let { history, children, routes, parseQueryString, stringifyQuery } = this.props;
    let createHistory = history ? () => history : createHashHistory;

    this.history = useRoutes(createHistory)({
      routes: createRoutes(routes || children),
      parseQueryString,
      stringifyQuery
    });

    this._unlisten = this.history.listen((error, state) => {
      if (error) {
        this.handleError(error);
      } else {
        //Rendering page after setting state, here to make sure the data is loaded `onUpdate` before we show the page
        if (this.state.firstRender) {
          this.setState({
            firstRender: false,
            ...state
          });
        } else {
          var _this = this;
          this.props.onUpdate(state, function () {
            _this.setState(state);
          });
        }
        //End change in `onUpdate` related change to get client side rendering to behave like now
      }
    });
  },
  //============================================================================
  componentWillUnmount() {
    if (this._unlisten) {
      this._unlisten();
    }
  },

  render() {
    let { location, routes, params, components } = this.state;
    let { createElement } = this.props;

    if (location == null) {
      return null; // Async match
    }

    const routingProps = {
      history: this.history,
      createElement,
      location,
      routes,
      params,
      components
    };

    return <RoutingContext { ...routingProps } />;
  }

});
/*eslint-enable*/

export default FluxibleRouter;
