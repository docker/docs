'use strict';

import { createSelector } from 'reselect';
import { Map } from 'immutable';
import {
  getTeamNameFromRouter,
  getOrgNameFromRouter
} from './router.js';
import { currentUserSelector } from './users';

/**
 * Returns the raw state.repositories map
 */
const getRawRepoState = state => state.repositories;

/**
 * repositories returns all repository records from
 * state.repositories.repositories.
 *
 * @return Immutable.Map  Map of RepositoryRecords
 */
export const repositories = state =>
  state.repositories.getIn(['repositories', 'entities', 'repo'], new Map());

export const sortedRepos = createSelector(
  repositories,
  currentUserSelector,
  (repos, user) => {

    // filter out all the ones from the current user's namespace
    let sortedArray = repos.toArray().filter((repo) => {
      return repo.namespace === user.name;
    });

    // now add back in all the ones from other namespaces
    return sortedArray.concat(repos.toArray().filter((repo) => {
      return repo.namespace !== user.name;
    }));
  }
);

/**
 * Returns the `repositoriesByNamespace` map. Used in:
 *  - organization list (to show first 10 repos for each org)
 *
 * @return Immutable.Map  Map of namespace names to normalized repository list
 * results
 */
export const getReposByNamespace = (state) =>
  state.repositories.get('repositoriesByNamespace', new Map());

/**
 *
 * @return Immutable.Map
 */
export const getTeamOrOrgRepos = createSelector(
  getTeamNameFromRouter,
  getOrgNameFromRouter,
  getRawRepoState,
  (teamName, orgName, repoState) => {
    if (teamName === undefined) {
      return repoState.getIn(['repositoriesByNamespace', orgName, 'entities', 'repo'], new Map());
    }
    return repoState.getIn(['repositoriesByTeam', orgName, teamName, 'entities', 'repo'], new Map());
  }
);

export const getTagsForRepo = (state, props) => state.repositories.getIn(['tags', props.params.namespace, props.params.name], new Map());
export const getAccessLevel = (state, props) => state.repositories.getIn(['userRepositoryAccess', window.user.name, props.params.namespace, props.params.name]);
export const getLongDescription = (state, props) => state.repositories.getIn(['repositories', props.params.namespace, props.params.name, 'repo', 'longDescription'], '');
export const getRepoForName = (state, props) => state.repositories.getIn(['repositories', props.params.namespace, props.params.name, 'repo'], new Map());

export const teamRepositoryAccessSelector = (state, props) => {
  return state.repositories.getIn(['repositories', props.params.namespace, props.params.name, 'access', 'teams'], new Map()).toJS();
};

export const getReposForUsername = (state, props) => state.repositories.getIn(['repositoriesByNamespace', props.params.username, 'entities', 'repo'], new Map());

// TODO: REMOVE ALL OF BELOW SELECTORS
export const namespaceSelector = () => window.user.name;
export const userRepositories = state => state.repositories.get('userRepositories', new Map());
export const flattenedSortedRepositoriesSelector = createSelector(
    [repositories],
    (repos) =>
      repos              // { namespace: { name: { repo, access, hidden: undefined }, ... }, ... }
        .map(x => x.toList())   // { namespace: [ { repo, access, hidden: undefined }, ... ],       ... }
        .toList()               // [ [ { repo, access, hidden: undefined }, ... ],       ... ]
        .flatten(1)             // [ { repo, access, hidden: undefined }, ... ]
        .sortBy((item) => -item.get('repo').id)
);
export const teamNamespaceAccessSelector = state => state.repositories.getIn(['teamNamespaceAccess']);
