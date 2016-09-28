'use strict';

import React, { Component, PropTypes } from 'react';
import { Map } from 'immutable';
const { array, func, object, number, string } = PropTypes;
import { connect } from 'react-redux';
import { reset as resetReduxForm } from 'redux-form';
import { replace } from 'react-router-redux';
import Spinner from 'components/common/spinner';
import RepositoryList from 'components/common/repositoryList';
import Button from 'components/common/button';
import consts from 'consts';
import CreateRepoForm from 'components/common/createRepoForm';
import { VelocityComponent } from 'velocity-react';
// selectors
import { createSelector, createStructuredSelector } from 'reselect';
import { sortedRepos } from 'selectors/repositories';
import { orgNamesSelector } from 'selectors/organizations';
import { currentUserSelector } from 'selectors/users';
import { getNamespaceNamesFromSearch } from 'selectors/namespaces';
// actions
import { mapActions } from 'utils';
import * as RepositoriesActions from 'actions/repositories';
import { listAdminOrUserOrganizations } from 'actions/organizations';
import { searchNamespaces } from 'actions/namespaces';
import autoaction from 'autoaction';
// misc
import uiDecorator from 'redux-ui';
import throttle from 'lodash/throttle';
import ZeroState from './zeroState.js';
import Filter from './filter.js';
import styles from './repositories.css';
import css from 'react-css-modules';

// This is used in getOrgsOrSearch selector to determine whether we have
// a search term
const getSearchTerm = (state) => state.ui.getIn(['app', 'repositories', 'searchTerm']);
const getFilter = (state) => state.ui.getIn(['app', 'repositories', 'filter']);

// getOrgsOrSearch returns the search results for namespaces if they
// exist, or a list of org names that the current user is a member of.
//
// This is because we want to show the org names that a user is a member of
// with nothing in the search bar (for quick access).
const getOrgsOrSearch = createSelector(
  [orgNamesSelector, getNamespaceNamesFromSearch, getSearchTerm, getFilter],
  (orgs, searchResults, term, filter) => {
    // This is a typeahead - users can search orgs by typing. We should show the
    // list of orgs if there's NO search term or the search term matches the
    // current fitler in the URL, as they've likely just clicked the dropdown
    if (term === '' || term === filter) {
      return orgs;
    }
    return searchResults;
  }
);

const repoState = state => state.repositories;

const allReposOrOrgRepos = createSelector(
  [getFilter, sortedRepos, repoState],
  (filter, sorted, repos) => {
    if (filter === '') {
      return sorted;
    }
    return repos.getIn(['repositoriesByNamespace', filter, 'entities', 'repo'], new Map()).toArray();
  }
);


const mapState = createStructuredSelector({
  router: state => state.router,
  repositories: (state) => state.repositories,
  orgs: getOrgsOrSearch,
  // This returns all org names to pass into CreateRepoForm
  orgNames: orgNamesSelector,
  repoList: allReposOrOrgRepos,
  user: currentUserSelector
});

@uiDecorator({
  // Selectors depend on this key being static
  key: 'repositories',
  state: {
    // filter represents the name of an organizationt to fitler repos by
    filter: (props, state) => { return state.router.locationBeforeTransitions.query.filter || ''; },
    // searchTerm represents the string a user has typed in the typeahead
    searchTerm: (props, state) => { return state.router.locationBeforeTransitions.query.filter || ''; },
    isCreateFormVisible: false
  },
  reducer: (state, action) => {
    // When we've successfully added a repository hide the form
    if (action.type === consts.repositories.CREATE_REPOSITORY && action.ready && !action.error) {
      return state.set('isCreateFormVisible', false);
    }
    return state;
  }
})
@connect(mapState, mapActions({
  listAdminOrUserOrganizations,
  ...RepositoriesActions,
  searchNamespaces,
  resetReduxForm,
  replace
}))
@autoaction({
  listRepositories: (props) => ({ namespace: props.ui.filter }),
  listAdminOrUserOrganizations: (props) => ({
    name: props.user.name,
    isAdmin: props.user.isAdmin
  })
}, { listAdminOrUserOrganizations, ...RepositoriesActions })
@css(styles)
export default class Repositories extends Component {

  static propTypes = {
    actions: object,
    currentPage: number,
    orgNames: array,
    orgs: array,
    user: object,
    repositories: object,
    repoList: array,
    filter: string,
    router: object,

    ui: object,
    updateUI: func.isRequired,
    location: object
  }
  static childContextTypes = {
      location: object
  }

  constructor(props, ctx, queue) {
    super(props, ctx, queue);

    this.searchNamespaces = throttle(
      ::this.props.actions.searchNamespaces,
      250
    );
  }

  getChildContext() {
      return {
          location: this.props.location
      };
  }

  selectNamespaceToFilter(_, { value: namespace }) {
    if (namespace === 'All accounts') {
      namespace = '';
    }

    // Get the current filter's repositories. An empty filter returns all.
    // Note that autoaction will load all data based on the UI prop change
    this.props.updateUI({
      filter: namespace,
      searchTerm: namespace
    });

    // Update the URL with a filter parmater if the filter exists
    const { pathname } = this.props.location;
    const filter = namespace ? '?filter=' + namespace : '';
    this.props.actions.replace(pathname + filter);
  }

  /**
   * When a user updates the namespace filter typeahead ensure we call the
   * search namespace action as a throttled function.
   */
  onTypeaheadKeydown(evt) {
    const { value } = evt.target;
    this.props.updateUI('searchTerm', value);
    this.searchNamespaces(value);
  }

  showCreateForm() {
    this.props.updateUI('isCreateFormVisible', true);
  }

  hideCreateForm() {
    this.props.updateUI('isCreateFormVisible', false);
  }

  createRepository(data) {
    const { namespace, name, shortDescription, visibility } = data;
    this.props.actions.createRepository({
      namespace,
      name,
      shortDescription,
      visibility
    });
  }

  renderCreateFormButton() {
    return (
      <span styleName='createFormButtonWrapper' id='new-repo-button'>
        <Button
          variant='secondary'
          onClick={ ::this.showCreateForm }
          disabled={ this.props.ui.isCreateFormVisible }>
          New repository
        </Button>
      </span>
    );
  }

  renderControls() {
    const { ui: { filter }, repoList } = this.props;
    const hideFilter = (filter === '' && repoList.length === 0);

    return (
      <div styleName='controls'>
        { hideFilter ? undefined :
            <Filter
              username={ this.props.user.name }
              orgs={ this.props.orgs }
              searchTerm={ this.props.ui.searchTerm }
              onChange={ ::this.onTypeaheadKeydown }
              onSelect={ ::this.selectNamespaceToFilter } /> }
        { this.renderCreateFormButton() }
      </div>
    );
  }

  renderZeroState() {
    const { isCreateFormVisible, filter } = this.props.ui;

    if (this.props.repoList && this.props.repoList.length) {
      return null;
    }

    // Don't show zero state if we have the new repo form open
    if (isCreateFormVisible) {
      return null;
    }

    if (filter) {
      return (
        <div styleName='zero'>
          <p>This account has no repositories.</p>
        </div>
      );
    }

    return <ZeroState />;
  }

  render() {
    const {
      ui: {
        isCreateFormVisible
      },
      user,
      orgNames
    } = this.props;

    let statuses = [
      [consts.repositories.LIST_REPOSITORIES, this.props.ui.filter]
    ];

    return (
      <div styleName='wrapper'>
        <div styleName='container'>
          { this.renderControls() }
          <VelocityComponent animation={ isCreateFormVisible ? 'slideDown' : 'slideUp' } duration={ 250 }>
            <div styleName='newFormWrapper'>
              <CreateRepoForm
                  onCancel={ ::this.hideCreateForm }
                  username={ user.name }
                  orgNames={ orgNames }
                  onSubmit={ ::this.createRepository } />
            </div>
          </VelocityComponent>
          <Spinner styleName='listContainer' loadingStatus={ statuses }>
            { this.renderZeroState() }
            { this.props.repoList && <RepositoryList
              repositories={ this.props.repoList } /> }
          </Spinner>
        </div>
      </div>
    );
  }
}
