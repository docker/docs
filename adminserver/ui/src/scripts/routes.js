'use strict';

import React from 'react';
import { IndexRedirect, Router, Route, IndexRoute, Redirect, browserHistory } from 'react-router';
import { syncHistoryWithStore, routerMiddleware } from 'react-router-redux';
import { createStore, applyMiddleware } from 'redux';
import { Provider } from 'react-redux';
import reducers from 'reducers';
import { middlewares } from 'middlewares';

import App from './app.js';

import {
  Login,
  Settings,
  API,
  Support,
  Repositories,
  Repository,
  RepositoryActivityTab,
  RepositoryDetailsTab,
  RepositoryTeamsTab,
  RepositoryUsersTab,
  RepositorySettingsTab,
  RepositoryTagsTab,
  GeneralSettings,
  StorageSettings,
  LicenseSettings,
  GarbageCollection,
  Organizations,
  Organization,
  Users,
  User,
  UserRepos,
  UserTeams,
  UserSettings,
  NotFound
} from 'components/scenes';

import OrgUsers from 'components/scenes/organization/users';
import OrgRepos from 'components/scenes/organization/repos';
import OrgSettings from 'components/scenes/organization/settings';

export const store = createStore(
    reducers,
    {},
    applyMiddleware(...middlewares)
);
const history = syncHistoryWithStore(browserHistory, store, {
    selectLocationState: (state) => state.router
});
applyMiddleware(routerMiddleware(history));

export const Root = (
  <Provider store={ store }>
    <Router history={ history }>
      <Route path='/' component={ App }>
        <IndexRedirect to='/repositories' />

        <Route path='orgs'>
          <IndexRoute component={ Organizations } />
          <Route path=':org' component={ Organization }>
            <IndexRedirect to='users' />
            <Route path='users' component={ OrgUsers } />
            <Route path='repos' component={ OrgRepos } />
            <Route path='settings' component={ OrgSettings } />
            <Route path='teams/:team'>
              <IndexRedirect to='users' />
              <Route path='users' component={ OrgUsers } />
              <Route path='repos' component={ OrgRepos } />
              <Route path='settings' component={ OrgSettings } />
            </Route>
          </Route>
        </Route>

        <Route path='users'>
          <IndexRoute component={ Users } />
          <Route path=':username' component={ User }>
            <IndexRedirect to='repos' />
            <Route path='repos' component={ UserRepos } />
            <Route path='teams' component={ UserTeams } />
            <Route path='settings' component={ UserSettings } />
          </Route>
        </Route>

        <Route path='repositories'>
          <IndexRoute component={ Repositories } />
          <Route path=':namespace/:name' component={ Repository } >
            <IndexRedirect to='details' />
            <Route path='details' component={ RepositoryDetailsTab } />
            <Route path='teams' component={ RepositoryTeamsTab } />
            <Route path='users' component={ RepositoryUsersTab } />
            <Route path='activity' component={ RepositoryActivityTab } />
            <Route path='settings' component={ RepositorySettingsTab } />
            <Route path='tags' component={ RepositoryTagsTab } />
          </Route>
        </Route>

        <Route path='support' component={ Support } />
        <Route path='docs/api' component={ API } />

        <Route path='admin/settings' component={ Settings }>
          <IndexRedirect to='general' />
          <Route path='general' component={ GeneralSettings } />
          <Route path='storage' component={ StorageSettings } />
          <Route path='license' component={ LicenseSettings } />
          <Route path='gc' component={ GarbageCollection } />
        </Route>

        <Route path='login' component={ Login } />

        <Redirect from='admin' to='/admin/settings/general' />
        <Redirect from='admin/settings' to='/admin/settings/general' />
        <Route path='*' component={ NotFound } />

      </Route>
    </Router>
  </Provider>
);
