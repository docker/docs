'use strict';

import createStore from 'fluxible/addons/createStore';
var debug = require('debug')('DashboardNamespacesStore');
import _ from 'lodash';

export default createStore({
  storeName: 'DashboardNamespacesStore',
  handlers: {
    RECEIVE_DASHBOARD_NAMESPACES: '_receiveOrgs',
    CURRENT_USER_CONTEXT: '_setContext',
    CREATE_REPO_RECEIVE_NAMESPACES: '_receiveOwnedNamespaces'
  },
  initialize() {
    this.namespaces = [];
    this.currentUserContext = '';
    this.ownedNamespaces = [];
  },
  _receiveOrgs(res) {
    //There are two API calls possible to get namespaces
    //`/v2/namespaces` -> returns `res.orgs.namespaces` an object with {namespaces: ['ns1', 'ns2', 'etc']}
    //`/v2/orgs` -> returns `res.orgs.results` with all orgs the user has read access on. We merge the `res.user` with this list
    if (res.orgs.namespaces) {
      this.namespaces = res.orgs.namespaces;
    } else if (res.orgs.results && _.isArray(res.orgs.results)) {
      var nsArray = _.pluck(res.orgs.results, 'orgname');
      nsArray.unshift(res.user);
      this.namespaces = nsArray;
    }
    this.emitChange();
  },
  _receiveOwnedNamespaces({
    namespaces, selectedNamespace
    }) {
    debug('receiving namespaces', namespaces, selectedNamespace);
    /**
     * namespaces is equivalent to the response in the namespaces API call
     */
    this.ownedNamespaces = namespaces.namespaces;
    this.emitChange();
  },
  _setContext: function({username}) {
    this.currentUserContext = username;
    this.emitChange();
  },
  getState() {
    return {
      currentUserContext: this.currentUserContext,
      namespaces: this.namespaces,
      ownedNamespaces: this.ownedNamespaces
    };
  },
  dehydrate() {
    return this.getState();
  },
  rehydrate(state) {
    this.currentUserContext = state.currentUserContext;
    this.namespaces = state.namespaces;
    this.ownedNamespaces = state.ownedNamespaces;
  }
});

