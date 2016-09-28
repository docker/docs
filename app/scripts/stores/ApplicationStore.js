'use strict';
const createStore = require('fluxible/addons/createStore');

var ApplicationStore = createStore({
  storeName: 'ApplicationStore',
  handlers: {
    CHANGE_ROUTE: 'handleNavigate'
  },
  initialize: function() {
    this.currentRoute = null;
  },
  handleNavigate: function(route) {
    if (this.currentRoute && route.path === this.currentRoute.path) {
      return;
    }

    this.currentRoute = route;
    this.emitChange();
  },
  getState: function() {
    return {
      route: this.currentRoute
    };
  },
  dehydrate: function() {
    return this.getState();
  },
  rehydrate: function(state) {
    this.currentRoute = state.route;
  }
});

module.exports = ApplicationStore;
