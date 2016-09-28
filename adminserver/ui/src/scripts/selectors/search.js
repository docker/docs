'use strict';

import { List } from 'immutable';
import { createSelector } from 'reselect';

const searchResults = (state) => state.search.get('search');

const searchResultCounts = createSelector(
  [searchResults],
  (results) => {
    return {
      orgs: results.getIn(['entities', 'org'], new List()).size,
      users: results.getIn(['entities', 'user'], new List()).size,
      repos: results.getIn(['entities', 'repo'], new List()).size
    };
  }
);

const searchResultPreview = createSelector(
  [searchResults],
  (results) => {
    let ret = {
      orgs: [],
      users: [],
      repos: []
    };
    results.getIn(['result'], new List()).forEach((result) => {
      const { id, schema } = result.toObject();
      // Schema is one of 'org', 'user', or 'repo'
      const sublist = ret[schema + 's'];
      if (sublist.length < 3) {
        sublist.push(results.getIn(['entities', schema, id]).toJS());
      }
    });
    return ret;
  }
);

// typeaheadPreview formats our search results into a list of objects;
// we create a a SearchResult for each entity type
const typeaheadPreview = createSelector(
  [searchResultPreview],
  (preview) => {
    return [
      ...Object.keys(preview.orgs).map((k, idx) => { return { data: preview.orgs[k], type: 'org', header: 'Organizations', index: idx }; }),
      ...Object.keys(preview.repos).map((k, idx) => { return { data: preview.repos[k], type: 'repo', header: 'Repositories', index: idx }; })
      // ...Object.keys(preview.users).map((k, idx) => { return { data: preview.users[k], type: 'user', header: 'Users', index: idx }; })
    ];
  }
);

export {
  searchResults,
  searchResultPreview,
  typeaheadPreview,
  searchResultCounts
};
