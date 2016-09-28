'use strict';

import { Map, List } from 'immutable';
import { createSelector } from 'reselect';

const allSearchResults = (state) => state.search.get('search', new Map());
const maxResults = (state) => state.ui.getIn(['app', 'globalSearch', 'maxResults']);
const maxResultsPerResource = (state) => state.ui.getIn(['app', 'globalSearch', 'maxResultsPerResource']);

const searchResultCounts = createSelector(
  [allSearchResults],
  (results) => {
    return {
      orgs: results.getIn(['entities', 'org'], new List()).size,
      users: results.getIn(['entities', 'user'], new List()).size,
      repos: results.getIn(['entities', 'repo'], new List()).size
    };
  }
);

/**
 * limitedSearchResults returns an array of search results in the same order as
 * the API response.
 *
 * Each result item has:
 *  - a .type field which lists whether the item is an org, user or repo.
 *  - a .data field with resource information
 *
 * The returned array will never exceed the maxResults UI state value in size.
 */
const limitedSearchResults = createSelector(
  [maxResults, maxResultsPerResource, allSearchResults],
  (max, maxPer, results) => {
    let result = [];

    const repos = results.getIn(['result', 'repos'], new List());
    const orgs = results.getIn(['result', 'orgs'], new List());
    const users = results.getIn(['result', 'users'], new List());

    // FUCK
    // TODO: results should be a single array containing an object of id/schema:
    // [
    //   {id: 'some-repo', schema: 'repo'},
    //   {id: 'some-org', schema: 'org'}
    // ]
    //
    // this allows us to maintain a strict order of the API search results, as
    // we can iterate over the result array in one sweep instead of separately.

    repos.forEach((id, i) => {
      if (i < maxPer && result.length < max) {
        result.push({
          data: results.getIn(['entities', 'repo', id]).toJS(),
          type: 'repo'
        });
      }
    });

    orgs.forEach((id, i) => {
      if (i < maxPer && result.length < max) {
        result.push({
          data: results.getIn(['entities', 'org', id]).toJS(),
          type: 'org'
        });
      }
    });

    users.forEach((id, i) => {
      if (i < maxPer && result.length < max) {
        result.push({
          data: results.getIn(['entities', 'user', id]).toJS(),
          type: 'user'
        });
      }
    });

    return result;
  }
);

/**
 * groupedLimitedSearchResults filters down allSearchResults into an array of entities
 * (orgs, users and repos). Each array will never be longer than the UI state's
 * maxResultsPerResource value.
 */
const groupedLimitedSearchResults = createSelector(
  [maxResultsPerResource, allSearchResults],
  (maxResource, results) => {
    let ret = {
      orgs: [],
      users: [],
      repos: []
    };

    // Iterate through the list to order results as returned from the search
    // API.
    results.getIn(['result'], new List()).forEach((result) => {
      const { id, schema } = result.toObject();
      // Schema is one of 'org', 'user', or 'repo'
      const sublist = ret[schema + 's'];
      if (sublist.length < maxResource) {
        sublist.push(results.getIn(['entities', schema, id]).toJS());
      }
    });

    return ret;
  }
);

/**
 * groupedSearchResults returns an array of search results where each resource
 * has the specific 'type' field added.
 *
 * This array will never be longer than the UI state's maxResults value
 */
const groupedSearchResults = createSelector(
  [maxResults, groupedLimitedSearchResults],
  (max, res) => {
    return [
      ...Object.keys(res.orgs).map((k, idx) => ({ data: res.orgs[k], type: 'org', header: 'Organizations', index: idx })),
      ...Object.keys(res.repos).map((k, idx) => ({ data: res.repos[k], type: 'repo', header: 'Repositories', index: idx }))
      // ...Object.keys(res.users).map((k, idx) => { return { data: res.users[k], type: 'user', header: 'Users', index: idx }; })
    ].slice(0, max);
  }
);

export {
  groupedSearchResults,
  limitedSearchResults,
  searchResultCounts
};
